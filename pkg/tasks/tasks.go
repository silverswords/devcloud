package tasks

import (
	"context"
	"io"
	"log"
	"os/exec"
	"strings"
)

// todo: read config and run tasks
func tripSpaceToCommandList(s string) []string {
	return strings.Split(s, " ")
}

func DockerRunMongo(ctx context.Context, req string) (string, error) {
	ExcuteCommand("docker rm -f mongo4")

	ExcuteCommand("docker run -d --name mongo4  -p 27017:27017 mongo:4")

	return "", nil
}

func ExcuteCommand(cmd string) {
	cmds := strings.Fields(cmd)

	args := []string{}
	if len(cmds) > 1 {
		args = cmds[1:]
	}
	process := exec.Command(cmds[0], args...)
	log.Printf("run command %+v", cmds)

	stdout, err := process.StdoutPipe()
	if err != nil {
		log.Fatal(cmd, err)
	}

	if err = process.Start(); err != nil {
		log.Fatal(cmd, err)
	}

	done := make(chan struct{}, 1)

	go func() {
		outputs := make([]byte, 512)
		line := make([]byte, 80)

		for {
			if n, err := stdout.Read(line); err != nil {
				if err == io.EOF {
					if n > 0 {
						outputs = append(outputs, line...)
					}
					break
				}

				log.Fatal(cmd, err)
			}

			outputs = append(outputs, line...)

			line = line[0:]
		}

		log.Print(string(outputs))

		done <- struct{}{}
	}()

	if err = process.Wait(); err != nil {
		log.Fatal(cmd, err)
	}

	<-done
}
