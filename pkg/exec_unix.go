//go:build !windows
// +build !windows

package qsnetcat

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

func ExecCommand(comm string, conn *QuantumSocket, interactive bool) error {
	params := strings.Split(comm, " ")
	cmd := &exec.Cmd{}
	defer conn.Close()

	if len(params) == 0 {
		return errors.New("no command specified")
	} else if len(params) == 1 {
		cmd = exec.Command(params[0])
	} else {
		cmd = exec.Command(params[0], params[1:]...)
	}
	os.Setenv("HISTFILE", "/dev/null")
	cmd.Env = append(cmd.Env, os.Environ()...)

	if interactive {
		// Start the command with a pty.
		ptmx, err := pty.Start(cmd)
		if err != nil {
			return err
		}

		// Make sure to close the pty at the end.
		defer ptmx.Close() // Best effort.

		// Handle pty size.
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		go func() {
			for range ch {
				if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
					logrus.Errorf("error resizing pty: %s", err)
				}
			}
		}()
		ch <- syscall.SIGWINCH // Initial resize.
		defer signal.Stop(ch)  // Cleanup signals when done.
		defer close(ch)

		// Set stdin in raw mode.
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState) // Best effort.

		// Copy stdin to the pty and the pty to stdout.
		// NOTE: The goroutine will keep reading until the next keystroke before returning.
		go func() { io.Copy(ptmx, conn) }()
		io.Copy(conn, ptmx)
		return nil
	}

	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	return cmd.Run()
}
