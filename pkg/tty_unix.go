//go:build aix || linux || solaris || zos
// +build aix linux solaris zos

package qsutils

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"golang.org/x/sys/unix"
)

var termState *unix.Termios = nil

func EnableInteractiveTerminal() error {
	oldState, err := MakeRawTerminal(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	termState = oldState
	return nil
}

func DisableInteractiveTerminal() {
	if termState != nil {
		SetCurrentTerminal(int(os.Stdin.Fd()), termState)
	}
}

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

	cmd := exec.Command("stty", "rows", fmt.Sprintf("%d", size[1]), "columns", fmt.Sprintf("%d", size[0]))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GetCurrentTerminal(fd int) (*unix.Termios, error) {
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}

	return termios, nil
}

func SetCurrentTerminal(fd int, term *unix.Termios) error {
	return unix.IoctlSetTermios(fd, unix.TCSETS, term)
}

func MakeRawTerminal(fd int) (*unix.Termios, error) {
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}

	oldState := *termios
	// This attempts to replicate the behaviour documented for cfmakeraw in
	// the termios(3) manpage.
	// termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	termios.Iflag |= unix.IGNPAR
	termios.Iflag &^= /* | unix.BRKINT | unix.PARMRK*/ unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON | unix.IEXTEN
	// termios.Oflag &^= unix.OPOST // We need this!!!!
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.ECHOE | unix.ECHOK
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	return &oldState, unix.IoctlSetTermios(fd, unix.TCSETS, termios)
}

func GetTerminalSize(fd int) (int, int, error) {
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return -1, -1, err
	}
	return int(ws.Col), int(ws.Row), nil
}
