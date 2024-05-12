//go:build windows
// +build windows

package qsnetcat

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/google/shlex"
	conpty "github.com/qsocket/conpty-go"
	"github.com/qsocket/qs-netcat/log"
	qsocket "github.com/qsocket/qsocket-go"
)

// const ATTACH_PARENT_PROCESS = ^uintptr(0)
const SHELL = "cmd.exe"

// func init() {
// 	proc := syscall.MustLoadDLL("KERNEL32.dll").MustFindProc("AttachConsole")
// 	proc.Call(ATTACH_PARENT_PROCESS) // We need this to get console output when using windowsgui subsystem.
// }

func ExecCommand(conn *qsocket.QSocket, specs *SessionSpecs) error {
	// If non specified spawn OS shell...
	if specs.Command == "" {
		specs.Command = SHELL
	}

	defer conn.Close()
	params, err := shlex.Split(specs.Command)
	if err != nil {
		return err
	}
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
	if specs.Interactive {
		cpty, err := conpty.Start(specs.Command)
		if err != nil {
			return err
		}
		defer cpty.Close()
		err = cpty.Resize(int(specs.TermSize.Cols), int(specs.TermSize.Rows))
		if err != nil {
			log.Error(err)
		}

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

func GetCurrentTermSize() (*Winsize, error) {
	return &Winsize{Cols: 80, Rows: 40}, nil
}
