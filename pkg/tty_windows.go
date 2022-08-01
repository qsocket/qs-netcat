//go:build windows
// +build windows

package qsutils

import (
	"fmt"
	"os"
	"time"
)

func SendTerminalSize(conn *QuantumSocket) error {
	w, h, err := GetTerminalSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	conn.SetWriteDeadline(time.Now().Add(time.Second * KNOCK_CHECK_DURATION))
	n, err := conn.Write([]byte{byte(w), byte(h)})
	conn.SetWriteDeadline(time.Time{})
	if err != nil {
		return err
	}
	if n != 2 {
		return ErrTtyFailed
	}

	return nil
}

func RecvTerminalSize(conn *QuantumSocket) error {
	size := make([]byte, 2)
	conn.SetReadDeadline(time.Now().Add(time.Second * KNOCK_CHECK_DURATION))
	n, err := conn.Read(size)
	conn.SetReadDeadline(time.Time{})
	if err != nil {
		return err
	}
	if n != 2 {
		return ErrTtyFailed
	}
	_, err = os.Stdin.WriteString(fmt.Sprintf("mode con: lines=%d cols=%d\r", size[1], size[0]))
	return err
}

func EnableInteractiveTerminal() error {
	// No such think in windows
	return nil
}

func DisableInteractiveTerminal() {
	// Done
}

func GetTerminalSize(fd int) (int, int, error) {
	return 80, 24, nil
}
