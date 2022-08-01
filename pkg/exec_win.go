//go:build windows
// +build windows

package qsutils

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func ExecCommand(comm string, conn *QuantumSocket) error {
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
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "HISTFILE=/dev/null")          // Disable history
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // Hide new process window
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	return cmd.Run()
}
