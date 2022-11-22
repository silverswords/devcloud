package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

const (
	windowsOS = "windows"
)

// installBinary installs files inside the dir directory.
func InstallBinary(version, binaryFilePrefix, githubRepo, dir string) error {
	var (
		err      error
		filepath string
	)

	filepath, err = downloadBinary(dir, version, binaryFilePrefix, githubRepo)
	if err != nil {
		return fmt.Errorf("error downloading %s binary: %w", binaryFilePrefix, err)
	}

	extractedFilePath, err := extractFile(filepath, dir, binaryFilePrefix)
	if err != nil {
		return err
	}

	// remove downloaded archive from the default disc bin path.
	err = os.Remove(filepath)
	if err != nil {
		return fmt.Errorf("failed to remove archive: %w", err)
	}

	binaryPath, err := moveFileToPath(extractedFilePath, dir)
	if err != nil {
		return fmt.Errorf("error moving %s binary to path: %w", binaryFilePrefix, err)
	}

	err = makeExecutable(binaryPath)
	if err != nil {
		return fmt.Errorf("error making %s binary executable: %w", binaryFilePrefix, err)
	}

	return nil
}

func makeExecutable(filepath string) error {
	if runtime.GOOS != windowsOS {
		err := os.Chmod(filepath, 0o777)
		if err != nil {
			return err
		}
	}

	return nil
}
func extractFile(filepath, dir, binaryFilePrefix string) (string, error) {
	var extractFunc func(string, string, string) (string, error)
	if archiveExt() == "zip" {
		extractFunc = unzipExternalFile
	} else {
		extractFunc = untarExternalFile
	}

	extractedFilePath, err := extractFunc(filepath, dir, binaryFilePrefix)
	if err != nil {
		return "", fmt.Errorf("error extracting %s binary: %w", binaryFilePrefix, err)
	}

	return extractedFilePath, nil
}

func unzipExternalFile(filepath, dir, binaryFilePrefix string) (string, error) {
	r, err := zip.OpenReader(filepath)
	if err != nil {
		return "", fmt.Errorf("error open zip file %s: %w", filepath, err)
	}
	defer r.Close()

	return unzip(&r.Reader, dir, binaryFilePrefix)
}

func unzip(r *zip.Reader, targetDir string, binaryFilePrefix string) (string, error) {
	foundBinary := ""
	for _, f := range r.File {
		fpath, err := sanitizeExtractPath(targetDir, f.Name)
		if err != nil {
			return "", err
		}

		if strings.HasSuffix(fpath, fmt.Sprintf("%s.exe", binaryFilePrefix)) {
			foundBinary = fpath
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		// #nosec G110
		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return "", err
		}
	}
	return foundBinary, nil
}

func untarExternalFile(filepath, dir, binaryFilePrefix string) (string, error) {
	reader, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("error open tar gz file %s: %w", filepath, err)
	}
	defer reader.Close()

	return untar(reader, dir, binaryFilePrefix)
}

func untar(reader io.Reader, targetDir string, binaryFilePrefix string) (string, error) {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	foundBinary := ""
	for {
		header, err := tr.Next()
		//nolint
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		} else if header == nil {
			continue
		}

		// untar all files in archive.
		path, err := sanitizeExtractPath(targetDir, header.Name)
		if err != nil {
			return "", err
		}

		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return "", err
			}
			continue
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
		if err != nil {
			return "", err
		}
		defer f.Close()

		// #nosec G110
		if _, err = io.Copy(f, tr); err != nil {
			return "", err
		}

		// If the found file is the binary that we want to find, save it and return later.
		if strings.HasSuffix(header.Name, binaryFilePrefix) {
			foundBinary = path
		}
	}
	return foundBinary, nil
}

// https://github.com/snyk/zip-slip-vulnerability, fixes gosec G305
func sanitizeExtractPath(destination string, filePath string) (string, error) {
	destpath := filepath.Join(destination, filePath)
	if !strings.HasPrefix(destpath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return "", fmt.Errorf("%s: illegal file path", filePath)
	}
	return destpath, nil
}
func moveFileToPath(path string, installLocation string) (string, error) {
	fileName := filepath.Base(path)
	destFilePath := ""

	destDir := installLocation
	destFilePath = filepath.Join(destDir, fileName)

	input, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	err = CreateDirectory(destDir)
	if err != nil {
		return "", err
	}

	// #nosec G306
	if err = os.WriteFile(destFilePath, input, 0o644); err != nil {
		if runtime.GOOS != windowsOS && strings.Contains(err.Error(), "permission denied") {
			err = errors.New(err.Error() + " - please run with sudo")
		}
		return "", err
	}

	if runtime.GOOS == windowsOS {
		p := os.Getenv("PATH")

		if !strings.Contains(strings.ToLower(p), strings.ToLower(destDir)) {
			pathCmd := "[System.Environment]::SetEnvironmentVariable('Path',[System.Environment]::GetEnvironmentVariable('Path','user') + '" + fmt.Sprintf(";%s", destDir) + "', 'user')"
			_, err := RunCmdAndWait("powershell", pathCmd)
			if err != nil {
				return "", err
			}
		}

		return fmt.Sprintf("%s\\discd.exe", destDir), nil
	}

	if strings.HasPrefix(fileName, "discd") && installLocation != "" {
		color.Set(color.FgYellow)
		fmt.Printf("\nDisc runtime installed to %s, you may run the following to add it to your path if you want to run discd directly:\n", destDir)
		fmt.Printf("    export PATH=$PATH:%s\n", destDir)
		color.Unset()
	}

	return destFilePath, nil
}
func archiveExt() string {
	ext := "tar.gz"
	if runtime.GOOS == windowsOS {
		ext = "zip"
	}

	return ext
}
func downloadBinary(dir, version, binaryFilePrefix, githubRepo string) (string, error) {
	fileURL := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/v%s/%s",

		githubRepo,
		version)

	return downloadFile(dir, fileURL)
}

func downloadFile(dir string, url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	filepath := path.Join(dir, fileName)
	_, err := os.Stat(filepath)
	if os.IsExist(err) {
		return "", nil
	}
	client := http.Client{ //nolint:exhaustruct
		Timeout: 0,
		Transport: &http.Transport{ //nolint:exhaustruct
			Dial: (&net.Dialer{ //nolint:exhaustruct
				Timeout: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   15 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			Proxy:                 http.ProxyFromEnvironment,
		},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("version not found from url: %s", url)
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = copyWithTimeout(context.Background(), out, resp.Body)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

/*
!
See: https://github.com/microsoft/vscode-winsta11er/blob/4b42060da64aea6f47adebe1dd654980ed87a046/common/common.go
Copyright (c) Microsoft Corporation. All rights reserved. Licensed under the MIT License.
*/
func copyWithTimeout(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	// Every 5 seconds, ensure at least 200 bytes (40 bytes/second average) are read.
	interval := 5
	minCopyBytes := int64(200)
	prevWritten := int64(0)
	written := int64(0)

	done := make(chan error)
	mu := sync.Mutex{}
	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer t.Stop()

	// Read the stream, 32KB at a time.
	go func() {
		var (
			writeErr, readErr     error
			writeBytes, readBytes int
			buf                   = make([]byte, 32<<10)
		)
		for {
			readBytes, readErr = src.Read(buf)
			if readBytes > 0 {
				// Write to disk and update the number of bytes written.
				writeBytes, writeErr = dst.Write(buf[0:readBytes])
				mu.Lock()
				written += int64(writeBytes)
				mu.Unlock()
				if writeErr != nil {
					done <- writeErr
					return
				}
			}
			if readErr != nil {
				// If error is EOF, means we read the entire file, so don't consider that as error.
				if !errors.Is(readErr, io.EOF) {
					done <- readErr
					return
				}

				// No error.
				done <- nil
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return written, ctx.Err()
		case <-t.C:
			mu.Lock()
			if written < prevWritten+minCopyBytes {
				mu.Unlock()
				return written, fmt.Errorf("stream stalled: received %d bytes over the last %d seconds", written, interval)
			}
			prevWritten = written
			mu.Unlock()
		case e := <-done:
			return written, e
		}
	}
}
