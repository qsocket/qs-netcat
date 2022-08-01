package qsutils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"qsutils/config"
	"qsutils/utils"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

const (
	// Qsocket options
	QSRN_GATE          = "qsocket.io"
	QSRN_GATE_TLS_PORT = 4443
	QSRN_GATE_PORT     = 80

	// OS Spesific binaries
	WIN_SHELL = "cmd.exe"
	NIX_SHELL = "/bin/bash -il"
)

var (
	ErrQsocketSessionEnd = errors.New("Qsocket session has ended")
	ErrTtyFailed         = errors.New("TTY initialization failed")
	spn                  = spinner.New(spinner.CharSets[9], 50*time.Millisecond)
)

func StartProbingQSRN(opts *config.Options) {
	var (
		qsrnAddr   = fmt.Sprintf("%s:%d", QSRN_GATE, QSRN_GATE_TLS_PORT)
		conn       any
		err        error
		firstProbe = true
	)

	if opts.DisableTLS {
		qsrnAddr = fmt.Sprintf("%s:%d", QSRN_GATE, QSRN_GATE_PORT)
	}

	go utils.WaitForExitSignal(os.Interrupt)

	for {
		if !firstProbe {
			time.Sleep(time.Duration(opts.ProbeInterval) * time.Second)
		} else {
			firstProbe = false
		}

		if opts.DisableTLS {
			conn, err = Dial(qsrnAddr, opts.UseTor)
			if err != nil {
				logrus.Error(err)
				continue
			}
		} else {
			conn, err = DialTLS(qsrnAddr, opts.UseTor)
			if err != nil {
				logrus.Error(err)
				continue
			}
		}

		qs, err := NewSocket(conn)
		if err != nil {
			logrus.Error(err)
			continue
		}

		err = SendKnockSequence(qs, opts.Secret, TagPortUsage(opts.Port))
		if err != nil {
			if err != ErrConnRefused && err != io.EOF {
				logrus.Error(err)
			}
			continue
		}

		// First check if forwarding enabled
		if opts.ForwardAddr != "" {
			// Redirect traffic to forward addr
			err = CreateOnConnectPipe(qs, fmt.Sprintf("%s:%d", opts.ForwardAddr, opts.Port))
			if err != nil {
				logrus.Error(err)
				continue
			}
		}

		// If non specified spawn OS shell...
		if opts.Execute == "" {
			opts.Execute = GetOsShell()
		}

		// If interactive, read terminal size first...
		if opts.Interactive {
			err = RecvTerminalSize(qs)
			if err != nil {
				logrus.Debug(err)
				continue
			}
		}

		// Execute command/program and redirect stdin/out/err
		err = ExecCommand(opts.Execute, qs)
		if err != nil && !strings.Contains(err.Error(), "connection reset by peer") {
			logrus.Error(err)
			continue
		}

	}
}

func BindSockets(con1, con2 *QuantumSocket) error {
	defer con1.Close()
	defer con2.Close()
	chan1 := CreateSocketChan(con1)
	chan2 := CreateSocketChan(con2)
	var err error
	for {
		select {
		case b1 := <-chan1:
			if b1 != nil {
				_, err = con2.Write(b1)
			} else {
				err = ErrQsocketSessionEnd
			}
		case b2 := <-chan2:
			if b2 != nil {
				_, err = con1.Write(b2)
			} else {
				err = ErrQsocketSessionEnd
			}
		}
		if err != nil {
			break
		}
	}
	return err
}

func CreateOnConnectPipe(con1 *QuantumSocket, addr string) error {
	defer con1.Close()
	chan1 := CreateSocketChan(con1)
	first := <-chan1
	if first == nil {
		return nil
	}

	logrus.Debug("Relaying first bytes!!")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	con2, err := NewSocket(conn)
	if err != nil {
		return err
	}
	defer con2.Close()
	_, err = con2.Write(first)
	if err != nil {
		return err
	}

	chan2 := CreateSocketChan(con2)

	for {
		select {
		case b1 := <-chan1:
			if b1 != nil {
				_, err = con2.Write(b1)
			} else {
				err = ErrQsocketSessionEnd
			}
		case b2 := <-chan2:
			if b2 != nil {
				_, err = con1.Write(b2)
			} else {
				err = ErrQsocketSessionEnd
			}
		}
		if err != nil {
			break
		}
	}
	return err
}

func ServeToLocal(opts *config.Options) {

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.Port))
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		inConn, err := ln.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		qs, err := NewSocket(inConn)
		if err != nil {
			logrus.Error(err)
			continue
		}
		ConnectAndBind(opts, qs)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func Connect(opts *config.Options) error {
	defer spn.Stop()

	if !opts.Quiet {
		spn.Suffix = " Dialing qsocket relay network..."
		spn.Start()
	}

	qsrnAddr := fmt.Sprintf("%s:%d", QSRN_GATE, QSRN_GATE_TLS_PORT)
	if opts.DisableTLS {
		qsrnAddr = fmt.Sprintf("%s:%d", QSRN_GATE, QSRN_GATE_PORT)
	}

	var (
		conn any
		err  error
	)

	if opts.DisableTLS {
		conn, err = Dial(qsrnAddr, opts.UseTor)
		if err != nil {
			return err
		}
	} else {
		conn, err = DialTLS(qsrnAddr, opts.UseTor)
		if err != nil {
			return err
		}
	}

	qs, err := NewSocket(conn)
	if err != nil {
		return err
	}

	err = SendKnockSequence(qs, opts.Secret, TagPortUsage(opts.Port))
	if err != nil {
		return err
	}

	return AttachToSocket(qs, opts.Interactive)
}

func ConnectAndBind(opts *config.Options, inConn *QuantumSocket) error {
	qsrnAddr := fmt.Sprintf("%s:%d", QSRN_GATE, QSRN_GATE_TLS_PORT)
	if opts.DisableTLS {
		qsrnAddr = fmt.Sprintf("%s:%d", QSRN_GATE, QSRN_GATE_PORT)
	}

	var (
		conn any
		err  error
	)

	if opts.DisableTLS {
		conn, err = Dial(qsrnAddr, opts.UseTor)
		if err != nil {
			return err
		}
	} else {
		conn, err = DialTLS(qsrnAddr, opts.UseTor)
		if err != nil {
			return err
		}
	}

	qs, err := NewSocket(conn)
	if err != nil {
		return err
	}

	err = SendKnockSequence(qs, opts.Secret, TagPortUsage(opts.Port))
	if err != nil {
		return err
	}

	return BindSockets(qs, inConn)
}

func AttachToSocket(conn *QuantumSocket, interactive bool) error {

	var err error
	if interactive {
		spn.Suffix = " Setting up TTY terminal..."
		err = SendTerminalSize(conn)
		if err != nil {
			return err
		}
		EnableInteractiveTerminal()
		defer DisableInteractiveTerminal()
	}
	spn.Stop()
	go func() {
		for {
			logrus.Debug("Reading from stdin...")
			n, readErr := io.Copy(conn, os.Stdin)
			if readErr != nil {
				err = readErr
				return
			}
			if n == 0 {
				logrus.Warn(ErrQsocketSessionEnd)
				break
			}
		}
	}()

	for {
		logrus.Debug("Reading from socket...")
		//_, err = writer2.ReadFrom(conn)
		n, writeErr := io.Copy(os.Stdout, conn)
		if writeErr != nil {
			err = writeErr
			break
		}
		if n == 0 {
			logrus.Warn(ErrQsocketSessionEnd)
			break
		}
	}
	return err
}

func Dial(addr string, tor bool) (net.Conn, error) {
	if tor {
		proxyDialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil,
			&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			},
		)
		if err != nil {
			return nil, err
		}

		conn, err := proxyDialer.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		return conn, nil

	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func DialTLS(addr string, tor bool) (net.Conn, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	if tor {
		proxyDialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil,
			&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			},
		)
		if err != nil {
			return nil, err
		}
		conn, err := proxyDialer.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		return tls.Client(conn, tlsConfig), nil

	}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func TagPortUsage(port int) uint8 {
	return uint8(utils.CaclChecksum([]byte{byte(port)}, 0xFF))
}

func GetOsShell() string {
	switch runtime.GOOS {
	case "linux", "darwin", "android", "freebsd", "ios", "netbsd", "openbsd", "solaris":
		return NIX_SHELL
	case "windows":
		return WIN_SHELL
	default:
		return NIX_SHELL
	}
}
