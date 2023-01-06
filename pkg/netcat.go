package qsnetcat

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/qsocket/qs-netcat/config"
	"github.com/qsocket/qs-netcat/utils"
	qsocket "github.com/qsocket/qsocket-go"

	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
	"golang.org/x/term"
)

var (
	ErrQsocketSessionEnd = errors.New("Qsocket session has ended")
	ErrTtyFailed         = errors.New("TTY initialization failed")
	ErrUntrustedCert     = errors.New("Certificate fingerprint mismatch")
	spn                  = spinner.New(spinner.CharSets[9], 50*time.Millisecond)
)

func StartProbingQSRN(opts *config.Options) {
	var err error
	go utils.WaitForExitSignal(os.Interrupt)
	// This is nessesary for persistence on windows
	SetWindowTitle(opts.Secret)
	time.Sleep(time.Duration(opts.ProbeInterval) * time.Second)
	for {
		qs := qsocket.NewSocket(opts.Secret, GetPeerTag(opts))
		if opts.UseTor {
			err = qs.DialProxy("socks5://127.0.0.1:9050")
		} else {
			err = qs.Dial(!opts.DisableTLS, opts.CertPinning)
		}
		if err != nil {
			if err != qsocket.ErrConnRefused {
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
			}
			continue
		}

		// If non specified spawn OS shell...
		if opts.Execute == "" {
			opts.Execute = OS_SHELL
		}

		// Execute command/program and redirect stdin/out/err
		err = ExecCommand(opts.Execute, qs, opts.Interactive)
		if err != nil && !strings.Contains(err.Error(), "connection reset by peer") {
			logrus.Error(err)
			continue
		}

	}
}

func SetWindowTitle(title string) {
	if runtime.GOOS == "windows" {
		exec.Command("cmd.exe", "/C", "title", title).Run()
	}
}

func CreateOnConnectPipe(qs *qsocket.Qsocket, addr string) error {
	defer qs.Close()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch := make(chan bool, 1)
	go func() {
		_, err = io.Copy(conn, qs)
		ch <- true
	}()
	go func() {
		_, err = io.Copy(qs, conn)
		ch <- true
	}()
	<-ch

	return err
}

func ServeToLocal(opts *config.Options) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.Port))
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		qs := qsocket.NewSocket(opts.Secret, GetPeerTag(opts))
		err = qs.Dial(!opts.DisableTLS, opts.CertPinning)
		if err != nil {
			logrus.Error(err)
			continue
		}
		go func() {
			_, err = io.Copy(conn, qs)
		}()
		_, err = io.Copy(qs, conn)
		if err != nil {
			logrus.Error(err)
		}
		qs.Close()
		conn.Close()
	}
}

func Connect(opts *config.Options) error {
	defer spn.Stop()
	if !opts.Quiet {
		spn.Suffix = " Dialing qsocket relay network..."
		spn.Start()
	}

	var err error
	qs := qsocket.NewSocket(opts.Secret, GetPeerTag(opts))
	if opts.UseTor {
		err = qs.DialProxy("socks5://127.0.0.1:9050")
	} else {
		err = qs.Dial(!opts.DisableTLS, opts.CertPinning)
	}
	if err != nil {
		return err
	}

	return AttachToSocket(qs, opts.Interactive)
}

func ConnectAndBind(opts *config.Options, inConn *qsocket.Qsocket) error {
	qs := qsocket.NewSocket(opts.Secret, GetPeerTag(opts))
	var err error
	if opts.UseTor {
		err = qs.DialProxy("socks5://127.0.0.1:9050")
	} else {
		err = qs.Dial(!opts.DisableTLS, opts.CertPinning)
	}
	if err != nil {
		return err
	}

	return qsocket.BindSockets(qs, inConn)
}

func AttachToSocket(conn *qsocket.Qsocket, interactive bool) error {
	defer conn.Close()
	if interactive {
		spn.Suffix = " Setting up TTY terminal..."
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

	spn.Stop()
	go func() { io.Copy(conn, os.Stdin) }()
	io.Copy(os.Stdout, conn)

	// go func() {
	// 	for {
	// 		logrus.Debug("Reading from stdin...")
	// 		n, readErr := io.Copy(conn, os.Stdin)
	// 		if readErr != nil {
	// 			err = readErr
	// 			logrus.Error("returning...")
	// 			return
	// 		}
	// 		if n == 0 {
	// 			logrus.Warn(ErrQsocketSessionEnd)
	// 			break
	// 		}
	// 	}
	// }()

	// for {
	// 	logrus.Debug("Reading from socket...")
	// 	//_, err = writer2.ReadFrom(conn)
	// 	n, writeErr := io.Copy(os.Stdout, conn)
	// 	if writeErr != nil {
	// 		logrus.Error("returning2...")
	// 		err = writeErr
	// 		break
	// 	}
	// 	if n == 0 {
	// 		logrus.Warn(ErrQsocketSessionEnd)
	// 		break
	// 	}
	// }
	return nil
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

func DialTLS(addr string, tor, certPinning bool) (net.Conn, error) {
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

	if certPinning {
		connState := conn.ConnectionState()
		for _, peerCert := range connState.PeerCertificates {
			hash := sha256.Sum256(peerCert.Raw)
			if !bytes.Equal(hash[0:], []byte(qsocket.CERT_FINGERPRINT)) {
				return nil, ErrUntrustedCert
			}
		}

	}

	return conn, nil
}

func GetPeerTag(opts *config.Options) byte {
	tag := byte(qsocket.TAG_PEER_CLI)
	if opts.Listen {
		tag = qsocket.TAG_PEER_SRV
	}

	if opts.Port != 0 ||
		opts.ForwardAddr != "" {
		return qsocket.TAG_PEER_PROXY
	}

	return tag
}
