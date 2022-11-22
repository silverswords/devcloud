package utils

import (
	"context"
	"os/exec"

	"github.com/docker/docker/client"
)

// IsDockerInstalled checks whether docker is installed/running.
func IsDockerInstalled() bool {
	//nolint:staticcheck
	cli, err := client.NewEnvClient()
	if err != nil {
		return false
	}
	_, err = cli.Ping(context.Background())
	return err == nil
}

func IsPodmanInstalled() bool {
	cmd := exec.Command("podman", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
