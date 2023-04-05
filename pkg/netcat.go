package qsnetcat

import (
	"errors"
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
	"golang.org/x/term"
)

var (
	ErrQsocketSessionEnd = errors.New("qsocket session has ended")
	ErrTtyFailed         = errors.New("TTY initialization failed")
	ErrUntrustedCert     = errors.New("certificate fingerprint mismatch")
	spn                  = spinner.New(spinner.CharSets[9], 50*time.Millisecond)
)

func StartProbingQSRN(opts *config.Options) {
	var (
		err      error
		firstRun bool = true
	)
	go utils.WaitForExitSignal(os.Interrupt)
	// This is nessesary for persistence on windows
	SetWindowTitle(opts.Secret) // This is nessesary for checking if the beacon is already running or not in Windows.
	os.Unsetenv("QS_ARGS")      // Remove this for allowing recursive qs-netcat usage

	for {
		if !firstRun {
			time.Sleep(time.Duration(opts.ProbeInterval) * time.Second)
		} else {
			firstRun = false
		}
		qs := qsocket.NewSocket(opts.Secret)
		err = qs.SetE2E(opts.End2End)
		if err != nil {
			logrus.Fatal(err)
		}
		err = qs.SetCertPinning(opts.CertPinning)
		if err != nil {
			logrus.Fatal(err)
		}
		err = qs.AddIdTag(GetPeerTag(opts))
		if err != nil {
			logrus.Fatal(err)
		}
		if opts.UseTor {
			err = qs.DialProxy("127.0.0.1:9050")
		} else {
			if opts.DisableEnc {
				err = qs.DialTCP()
			} else {
				err = qs.Dial()
			}
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
			err = CreateOnConnectPipe(qs, opts.ForwardAddr)
			if err != nil {
				logrus.Error(err)
			}
			continue
		}

		// If non specified spawn OS shell...
		if opts.Execute == "" {
			opts.Execute = SHELL
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

func CreateOnConnectPipe(qs *qsocket.QSocket, addr string) error {
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
	ln, err := net.Listen("tcp", opts.ForwardAddr)
	if err != nil {
		logrus.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		qs := qsocket.NewSocket(opts.Secret)
		err = qs.SetE2E(opts.End2End)
		if err != nil {
			logrus.Fatal(err)
		}
		err = qs.SetCertPinning(opts.CertPinning)
		if err != nil {
			logrus.Fatal(err)
		}
		err = qs.AddIdTag(GetPeerTag(opts))
		if err != nil {
			logrus.Fatal(err)
		}

		if opts.DisableEnc {
			err = qs.DialTCP()
		} else {
			err = qs.Dial()
		}
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
	if opts.ForwardAddr != "" {
		ServeToLocal(opts)
		return nil
	}
	defer spn.Stop()
	if !opts.Quiet {
		spn.Suffix = " Dialing qsocket relay network..."
		spn.Start()
	}
	var err error
	qs := qsocket.NewSocket(opts.Secret)
	err = qs.SetE2E(opts.End2End)
	if err != nil {
		logrus.Fatal(err)
	}
	err = qs.SetCertPinning(opts.CertPinning)
	if err != nil {
		logrus.Fatal(err)
	}
	err = qs.AddIdTag(GetPeerTag(opts))
	if err != nil {
		logrus.Fatal(err)
	}

	if opts.UseTor {
		err = qs.DialProxy("socks5://127.0.0.1:9050")
	} else {
		if opts.DisableEnc {
			err = qs.DialTCP()
		} else {
			err = qs.Dial()
		}
	}
	if err != nil {
		return err
	}
	return AttachToSocket(qs, opts.Interactive)
}

func AttachToSocket(conn *qsocket.QSocket, interactive bool) error {
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
	return nil
}

func GetPeerTag(opts *config.Options) byte {
	tag := byte(qsocket.TAG_PEER_CLI)
	if opts.Listen {
		tag = qsocket.TAG_PEER_SRV
	}
	if opts.ForwardAddr != "" {
		return qsocket.TAG_PEER_PROXY
	}
	return tag
}
