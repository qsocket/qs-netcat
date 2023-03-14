//go:build android
// +build android

package qsnetcat

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	qsocket "github.com/qsocket/qs-netcat/qsocket-go"
	"github.com/sirupsen/logrus"
	// _ "golang.org/x/mobile/app"
)

const SHELL = "sh"

func ExecCommand(comm string, conn *qsocket.QSocket, interactive bool) error {
	defer conn.Close()
	params := strings.Split(comm, " ")
	ncDir, err := os.Executable() // Get the full path of the executalbe.
	if err != nil {
		return err
	}
	os.Setenv("qs_netcat", ncDir)      // Set binary dir to env variable for easy access.
	os.Setenv("HISTFILE", "/dev/null") // Unset histfile for disabling logging.
	cmd := exec.Command(params[0])
	if len(params) > 1 {
		cmd = exec.Command(params[0], params[1:]...)
	}

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
