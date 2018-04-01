package docker

import (
	"os"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/term"
)

// Controller is a Docker implementation of container.Controller.
type Controller struct {
	client *client.Client

	stdin  termStream
	stdout termStream
	stderr termStream
}

type termStream struct {
	stream *os.File
	fd     uintptr
}

// NewController returns a `*Controller` with stdin, stdout and stderr initialized.
func NewController() (*Controller, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	stdinfd, _ := term.GetFdInfo(os.Stdin)
	stdoutfd, _ := term.GetFdInfo(os.Stdout)
	stderrfd, _ := term.GetFdInfo(os.Stderr)

	return &Controller{
		client: cli,
		stdin:  termStream{stream: os.Stdin, fd: stdinfd},
		stdout: termStream{stream: os.Stdout, fd: stdoutfd},
		stderr: termStream{stream: os.Stderr, fd: stderrfd},
	}, nil
}
