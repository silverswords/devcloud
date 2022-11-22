package utils

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

func CreateDirectory(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}
	return os.Mkdir(dir, 0o777)
}

func RunCmdAndWait(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	resp, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	errB, err := io.ReadAll(stderr)
	if err != nil {
		//nolint
		return "", nil
	}

	err = cmd.Wait()
	if err != nil {
		// in case of error, capture the exact message.
		if len(errB) > 0 {
			return "", errors.New(string(errB))
		}
		return "", err
	}

	return string(resp), nil
}
