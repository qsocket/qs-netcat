package qsnetcat

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/qsocket/qs-netcat/config"
	"github.com/qsocket/qs-netcat/log"
	qsocket "github.com/qsocket/qsocket-go"

	"github.com/briandowns/spinner"
	"golang.org/x/term"
)

var (
	ErrQsocketSessionEnd = errors.New("QSocket session has ended")
	ErrTtyFailed         = errors.New("TTY initialization failed")
	ErrUntrustedCert     = errors.New("Certificate fingerprint mismatch")
	spn                  = spinner.New(spinner.CharSets[9], 50*time.Millisecond)
)

func ProbeQSRN(opts *config.Options) error {
	// This is nessesary for persistence on windows
	os.Unsetenv("QS_ARGS") // Remove this for allowing recursive qs-netcat usage
	qs := qsocket.NewSocket(opts.Secret)
	err := qs.SetE2E(opts.End2End)
	if err != nil {
		return err
	}
	if opts.CertFingerprint != "" {
		err = qs.SetCertFingerprint(opts.CertFingerprint)
		if err != nil {
			return err
		}
	}
	err = qs.SetIdTag(GetPeerTag(opts))
	if err != nil {
		return err
	}
	if opts.SocksAddr != "" {
		err = qs.SetProxy(opts.SocksAddr)
		if err != nil {
			return err
		}
	}

	// Dial QSRN...
	if opts.DisableEnc {
		err = qs.DialTCP()
	} else {
		err = qs.Dial()
	}
	if err != nil {
		return err
	}

	log.Debug("Recving session specs...")
	specs, err := RecvSessionSpecs(qs, opts)
	if err != nil {
		return err
	}

	log.Info("Starting new session...")
	// Resize terminal with client dimentions...
	log.Debugf("Got term size: %dx%d", specs.TermSize.Rows, specs.TermSize.Cols)
	log.Debugf("Got command: %s", specs.Command)
	log.Debugf("Got forward: %s", specs.ForwardAddr)

	// First check if forwarding enabled
	if specs.ForwardAddr != "" {
		// Redirect traffic to forward addr
		err = CreateOnConnectPipe(qs, specs.ForwardAddr)
		if err != nil {
			return err
		}
	}

	go func() {
		// Execute command/program and redirect stdin/out/err
		err = ExecCommand(qs, specs)
		if err != nil {
			log.Error(err)
		}
	}()
	return err
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

func ServeToLocal(qs *qsocket.QSocket, opts *config.Options) {
	ln, err := net.Listen("tcp", ":"+strings.Split(opts.ForwardAddr, ":")[0])
	if err != nil {
		log.Fatal(err)
	}

	for {
		spn.Suffix = " Waiting for local connection..."
		spn.Start()
		conn, err := ln.Accept()
		if err != nil {
			log.Error(err)
			continue
		}
		spn.Suffix = " Dialing qsocket relay network..."
		if opts.SocksAddr != "" {
			err = qs.SetProxy(opts.SocksAddr)
			if err != nil {
				log.Fatal(err)
			}
		}
		if opts.DisableEnc {
			err = qs.DialTCP()
		} else {
			err = qs.Dial()
		}
		if err != nil {
			spn.Stop()
			log.Error(err)
			continue
		}
		spn.Suffix = " Forwarding local traffic..."
		go func() {
			_, err = io.Copy(conn, qs)
		}()
		_, err = io.Copy(qs, conn)
		if err != nil {
			log.Debug(err)
		}
		qs.Close()
		conn.Close()
	}
}

func Connect(opts *config.Options) error {
	defer spn.Stop()
	if opts == nil {
		return errors.New("options are not initialized")
	}
	if !opts.Quiet {
		spn.Suffix = " Dialing qsocket relay network..."
		spn.Start()
	}

	qs := qsocket.NewSocket(opts.Secret)
	err := qs.SetE2E(opts.End2End)
	if err != nil {
		log.Fatal(err)
	}
	if opts.CertFingerprint != "" {
		err = qs.SetCertFingerprint(opts.CertFingerprint)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = qs.SetIdTag(GetPeerTag(opts))
	if err != nil {
		log.Fatal(err)
	}

	if opts.ForwardAddr != "" {
		parts := strings.Split(opts.ForwardAddr, ":")
		switch len(parts) {
		case 2:
		case 3:
			opts.ForwardAddr = fmt.Sprintf("%s:%s", parts[1], parts[2])
			ServeToLocal(qs, opts)
		default:
			spn.Stop()
			log.Fatal("Invalid forward address!")
		}
	}
	if opts.SocksAddr != "" {
		err = qs.SetProxy(opts.SocksAddr)
		if err != nil {
			log.Fatal(err)
		}
	}
	if opts.DisableEnc {
		err = qs.DialTCP()
	} else {
		err = qs.Dial()
	}
	if err != nil {
		return err
	}

	log.Debug("Sending session specs...")
	err = SendSessionSpecs(qs, opts)
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
	if opts.Listen {
		return qsocket.PEER_SRV
	}
	return qsocket.PEER_CLI
}
