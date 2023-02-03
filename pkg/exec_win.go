//go:build windows
// +build windows

package qsnetcat

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	conpty "github.com/EgeBalci/conpty-go"
	qsocket "github.com/qsocket/qsocket-go"
)

const SHELL = "cmd.exe"

func ExecCommand(comm string, conn *qsocket.Qsocket, interactive bool) error {
	defer conn.Close()
	params := strings.Split(comm, " ")
	ncDir, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}
	os.Setenv("qs_netcat", ncDir)
	cmd := exec.Command(params[0])
	if len(params) > 1 {
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
