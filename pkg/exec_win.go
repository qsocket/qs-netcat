//go:build windows
// +build windows

package qsnetcat

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	conpty "github.com/EgeBalci/conpty-go"
	qsocket "github.com/qsocket/qsocket-go"
)

const OS_SHELL = "cmd.exe"

func ExecCommand(comm string, conn *qsocket.Qsocket, interactive bool) error {
	defer conn.Close()
	params := strings.Split(comm, " ")
	cmd := &exec.Cmd{Env: os.Environ()}
	defer conn.Close()
	if len(params) == 0 {
		return errors.New("no command specified")
	} else if len(params) == 1 {
		cmd = exec.Command(params[0])
	} else {
		cmd = exec.Command(params[0], params[1:]...)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // Hide new process window
	if interactive {
		cpty, err := conpty.Start(comm)
		if err != nil {
			return err
		}
		defer cpty.Close()

		go func() {
			go io.Copy(conn, cpty)
			io.Copy(cpty, conn)
		}()

		_, err = cpty.Wait(context.Background())
		return err
	}

	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	return cmd.Run()
}
