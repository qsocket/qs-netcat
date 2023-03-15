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

	conpty "github.com/qsocket/conpty-go"
	qsocket "github.com/qsocket/qsocket-go"
)

func init() {
	const ATTACH_PARENT_PROCESS = ^uintptr(0)
	proc := syscall.MustLoadDLL("KERNEL32.dll").MustFindProc("AttachConsole")
	proc.Call(ATTACH_PARENT_PROCESS) // We need this to get console output when using windowsgui subsystem.
}

const SHELL = "cmd.exe"

func ExecCommand(comm string, conn *qsocket.QSocket, interactive bool) error {
	defer conn.Close()
	params := strings.Split(comm, " ")
	ncDir, err := filepath.Abs(os.Args[0]) // Get the full path of the executalbe.
	if err != nil {
		return err
	}
	os.Setenv("qs_netcat", ncDir) // Set binary dir to env variable for easy access.
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
