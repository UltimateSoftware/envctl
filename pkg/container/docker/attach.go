package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	gosignal "os/signal"

	"github.com/UltimateSoftware/envctl/pkg/container"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/term"
)

// Attach attaches the terminal session of the currently running
// program to the container interactively.
func (c *Controller) Attach(m container.Metadata) error {
	restoreStdout, restoreStdin, err := c.makeRawTerminal()
	if err != nil {
		return err
	}

	acfg := types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	}

	resp, err := c.client.ContainerAttach(context.Background(), m.ID, acfg)
	defer resp.Close()
	if err != nil {
		return err
	}

	errchan := make(chan error)
	donechan := make(chan struct{})

	err = c.client.ContainerStart(
		context.Background(),
		m.ID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return err
	}

	c.mirrorContainerTTY(m.ID)

	go func() {
		_, err := io.Copy(c.stdout.stream, resp.Reader)
		if err != nil {
			errchan <- err
			restoreStdout()
			return
		}

		restoreStdout()
		donechan <- struct{}{}
	}()

	go func() {
		_, err := io.Copy(resp.Conn, c.stdin.stream)
		if err != nil {
			errchan <- err
			restoreStdin()
			return
		}

		restoreStdin()
		resp.CloseWrite()
	}()

	// Depending on the underlying image's entrypoint, there could be cases
	// where there's no command prompt. This could trick the user into thinking
	// that the process is hung, when in fact there just hasn't been anything
	// to write to stdout.
	fmt.Fprintf(
		c.stdout.stream,
		"If you don't see a command prompt, try pressing enter.\r\n",
	)

	select {
	case err = <-errchan:
		if err != nil {
			return err
		}
	case <-donechan:
	}

	return nil
}

// makeRawTerminal sets the terminal currently pointed to by stdin and stdout
// into a raw terminal. This is necessary for communication with the Docker
// container over the attach socket. If anything goes wrong, it returns an error
// and tries to restore the terminal to its previous state, but it might not
// succeed in doing so depending on what the issue was. If it was successful,
// it returns two callback functions for the caller to restore the terminal
// back to its previous state when ready. The first is for stdout, the second
// is for stderr.
func (c *Controller) makeRawTerminal() (func() error, func() error, error) {
	// This stuff is required to make interactive sessions in the container
	// less buggy. For example, without it, any command typed at the prompt will
	// get repeated out before printing the execution results.
	oldStdout, err := term.MakeRaw(c.stdout.fd)
	if err != nil {
		return nil, nil, err
	}

	restoreStdout := func() error {
		return term.RestoreTerminal(c.stdout.fd, oldStdout)
	}

	oldStdin, err := term.MakeRaw(c.stdin.fd)
	if err != nil {
		term.RestoreTerminal(c.stdout.fd, oldStdout)
		return nil, nil, err
	}

	restoreStdin := func() error {
		return term.RestoreTerminal(c.stdin.fd, oldStdin)
	}

	return restoreStdout, restoreStdin, nil
}

func (ts *termStream) getTTYSize() (uint, uint) {
	ws, err := term.GetWinsize(ts.fd)
	if err != nil {
		if ws == nil {
			return 0, 0
		}
	}
	return uint(ws.Width), uint(ws.Height)
}

// mirrorContainerTTY handles keeping the tty dimensions in sync from the host
// to the container.
func (c *Controller) mirrorContainerTTY(cntid string) error {
	handleTerminalResize := func() {
		width, height := c.stdout.getTTYSize()
		if width == 0 && height == 0 {
			return
		}

		options := types.ResizeOptions{
			Width:  width,
			Height: height,
		}

		c.client.ContainerResize(context.Background(), cntid, options)
	}

	// Run this the first time to establish the link between the container's TTY
	// and the terminal emulator's TTY.
	handleTerminalResize()

	sigchan := make(chan os.Signal, 1)
	gosignal.Notify(sigchan, signal.SIGWINCH)
	go func() {
		for range sigchan {
			handleTerminalResize()
		}
	}()
	return nil
}
