package qsocket

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"runtime"
	"time"

	stream "github.com/EgeBalci/encrypted-stream"
	"golang.org/x/net/proxy"
)

// Knock tags
// 000 000 0 0
// |   |   | |
// [ARCH]  | |
//
//	    |   | |
//	    [OS]| |
//		       | |
//		       [PROXY]
//		         [SRV|CLI]
const (
	// TAG_ARCH_AMD64 Tag ID value representing connections from devices with AMD64 architecture.
	TAG_ARCH_AMD64 = 0xE0 // 00110000 => AMD64
	// TAG_ARCH_386 Tag ID value representing connections from devices with 386 architecture.
	TAG_ARCH_386 = 0x20 // 00100000 => 386
	// TAG_ARCH_ARM64 Tag ID value representing connections from devices with ARM64 architecture.
	TAG_ARCH_ARM64 = 0x40 // 01000000 => ARM64
	// TAG_ARCH_ARM Tag ID value representing connections from devices with ARM architecture.
	TAG_ARCH_ARM = 0x60 // 01100000 => ARM
	// TAG_ARCH_MIPS64 Tag ID value representing connections from devices with MIPS64 architecture.
	TAG_ARCH_MIPS64 = 0x80 // 10000000 => MIPS64
	// TAG_ARCH_MIPS Tag ID value representing connections from devices with MIPS architecture.
	TAG_ARCH_MIPS = 0xA0 // 10100000 => MIPS
	// TAG_ARCH_MIPS64LE Tag ID value representing connections from devices with MIPS64LE architecture.
	TAG_ARCH_MIPS64LE = 0xC0 // 11000000 => MIPS64LE
	// TAG_ARCH_UNKNOWN Tag ID value representing connections from devices with UNKNOWN architecture.
	TAG_ARCH_UNKNOWN = 0x00 // 11100000 => UNKNOWN

	// TAG_OS_LINUX Tag ID value representing connections from LINUX devices.
	TAG_OS_LINUX = 0x1C // 00000000 => LINUX
	// TAG_OS_DARWIN Tag ID value representing connections from DARWIN devices.
	TAG_OS_DARWIN = 0x04 // 00000100 => DARWIN
	// TAG_OS_WINDOWS Tag ID value representing connections from WINDOWS devices.
	TAG_OS_WINDOWS = 0x08 // 00001000 => WINDOWS
	// TAG_OS_ANDROID Tag ID value representing connections from ANDROID devices.
	TAG_OS_ANDROID = 0x0C // 00001100 => ANDROID
	// TAG_OS_IOS Tag ID value representing connections from IOS devices.
	TAG_OS_IOS = 0x10 // 00010000 => IOS
	// TAG_OS_FREEBSD Tag ID value representing connections from FREEBSD devices.
	TAG_OS_FREEBSD = 0x14 // 00010100 => FREEBSD
	// TAG_OS_OPENBSD Tag ID value representing connections from OPENBSD devices.
	TAG_OS_OPENBSD = 0x18 // 00011000 => MIPS64LE
	// TAG_OS_UNKNOWN Tag ID value representing connections from UNKNOWN devices.
	TAG_OS_UNKNOWN = 0x00 // 00011100 => UNKNOWN

	// TAG_PEER_PROXY Tag ID for representing proxy mode connections.
	TAG_PEER_PROXY = 0x02 // 00000010 => Proxy connection
	// Tag ID for representing server mode connections.
	TAG_PEER_SRV = 0x00 // 00000000 => Server
	// Tag ID for representing client mode connections.
	TAG_PEER_CLI = 0x01 // 00000001 => Client
)

var (
	ErrUntrustedCert       = errors.New("certificate fingerprint mismatch")
	ErrUninitializedSocket = errors.New("socket not initiated")
	ErrQSocketSessionEnd   = errors.New("qSocket session has ended")
	ErrUnexpectedSocket    = errors.New("unexpected socket type")
	ErrInvalidIdTag        = errors.New("invalid peer ID tag")
	ErrNoTlsConnection     = errors.New("TLS socket is nil")
)

// A QSocket structure contains required values
// for performing a knock sequence with the QSRN gate.
//
// `Secret` value can be considered as the password for the QSocket connection,
// It will be used for generating a 128bit unique identifier (UID) for the connection.
//
// `tag` value is used internally for QoS purposes.
// It specifies the operating system, architecture and the type of connection initiated by the peers,
// the relay server uses these values for optimizing the connection performance.
type QSocket struct {
	Secret     string
	CertVerify bool
	tag        byte
	conn       net.Conn
	tlsConn    *tls.Conn
	encConn    *stream.EncryptedStream
}

// NewSocket creates a new QSocket structure with the given secret.
// `certVerify` value is used for enabling the certificate validation on TLS connections
func NewSocket(secret string, certVerify bool) *QSocket {
	tag := GetDefaultTag()
	return &QSocket{
		Secret:     secret,
		CertVerify: certVerify,
		tag:        tag,
		conn:       nil,
		tlsConn:    nil,
		encConn:    nil,
	}
}

// AddIdTag adds a peer identification tag to the QSocket.
func (qs *QSocket) AddIdTag(idTag byte) error {
	switch idTag {
	case TAG_PEER_SRV:
		qs.tag |= idTag
	case TAG_PEER_CLI:
		qs.tag |= idTag
	case TAG_PEER_PROXY:
		qs.tag |= idTag
	case TAG_PEER_PROXY | TAG_PEER_CLI:
		qs.tag |= idTag
	default:
		return ErrInvalidIdTag
	}
	return nil
}

// DialTCP creates a TCP connection to the `QSRN_GATE` on `QSRN_GATE_PORT`.
func (qs *QSocket) DialTCP() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", QSRN_GATE, QSRN_GATE_PORT))
	if err != nil {
		return err
	}
	qs.conn = conn
	return qs.SendKnockSequence()
}

// Dial creates a TLS connection to the `QSRN_GATE` on `QSRN_GATE_TLS_PORT`.
// Based on the `VerifyCert` parameter, certificate fingerprint validation (a.k.a. SSL pinning)
// will be performed after establishing the TLS connection.
func (qs *QSocket) Dial() error {
	conf := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", QSRN_GATE, QSRN_GATE_TLS_PORT), conf)
	if err != nil {
		return err
	}
	qs.tlsConn = conn

	if qs.CertVerify {
		connState := conn.ConnectionState()
		for _, peerCert := range connState.PeerCertificates {
			hash := sha256.Sum256(peerCert.Raw)
			if !bytes.Equal(hash[0:], []byte(CERT_FINGERPRINT)) {
				return ErrUntrustedCert
			}
		}
	}

	err = qs.SendKnockSequence()
	if err != nil {
		return err
	}

	return qs.InitE2E()

}

func (qs *QSocket) InitE2E() error {
	if qs.tlsConn == nil { // We need a valid TLS connection for initiating PAKE for E2E.
		return ErrNoTlsConnection
	}

	sum := sha256.Sum256([]byte(qs.Secret))
	cipher, err := stream.NewAESGCMCipher(sum[:])
	if err != nil {
		return err
	}

	config := &stream.Config{
		Cipher:                   cipher,
		DisableNonceVerification: true, // This is nessesary because we don't really know who (client/server) speaks first on the relay connection.
	}

	// Create an encrypted stream from a conn.
	encryptedConn, err := stream.NewEncryptedStream(qs.tlsConn, config)
	if err != nil {
		return err
	}
	qs.encConn = encryptedConn
	return nil
}

// DialProxy tries to create TCP connection to the `QSRN_GATE` using a SOCKS5 proxy.
// `proxyAddr` should contain a valid SOCKS5 proxy.
func (qs *QSocket) DialProxy(proxyAddr string) error {
	proxyDialer, err := proxy.SOCKS5("tcp", proxyAddr, nil,
		&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		},
	)
	if err != nil {
		return err
	}

	conn, err := proxyDialer.Dial("tcp", fmt.Sprintf("%s:%d", QSRN_GATE, QSRN_GATE_PORT))
	if err != nil {
		return err
	}
	qs.conn = conn
	return qs.SendKnockSequence()
}

// IsClosed checks if the QSocket connection to the `QSRN_GATE` is ended.
func (qs *QSocket) IsClosed() bool {
	return qs.conn == nil && qs.tlsConn == nil
}

// SetReadDeadline sets the read deadline on the underlying connection.
// A zero value for t means Read will not time out.
func (qs *QSocket) SetReadDeadline(t time.Time) error {
	if qs.tlsConn != nil {
		return qs.tlsConn.SetReadDeadline(t)
	}
	if qs.conn != nil {
		return qs.conn.SetReadDeadline(t)
	}
	return nil
}

// SetWriteDeadline sets the write deadline on the underlying connection.
// A zero value for t means Write will not time out.
// After a Write has timed out, the TLS state is corrupt and all future writes will return the same error.
// Even if write times out, it may return n > 0, indicating that some of the data was successfully written. A zero value for t means Write will not time out.
func (qs *QSocket) SetWriteDeadline(t time.Time) error {
	if qs.tlsConn != nil {
		return qs.tlsConn.SetWriteDeadline(t)
	}
	if qs.conn != nil {
		return qs.conn.SetWriteDeadline(t)
	}
	return nil
}

// RemoteAddr returns the remote network address.
func (qs *QSocket) RemoteAddr() net.Addr {
	if qs.tlsConn != nil {
		return qs.tlsConn.RemoteAddr()
	}
	if qs.conn != nil {
		return qs.conn.RemoteAddr()
	}
	return nil
}

// LocalAddr returns the local network address.
func (qs *QSocket) LocalAddr() net.Addr {
	if qs.tlsConn != nil {
		return qs.tlsConn.LocalAddr()
	}
	if qs.conn != nil {
		return qs.conn.LocalAddr()
	}
	return nil
}

// TLS checks if the underlying connection is TLS or not.
func (qs *QSocket) TLS() bool {
	return qs.tlsConn != nil
}

// Read reads data from the connection.
//
// As Read calls Handshake, in order to prevent indefinite blocking a deadline must be set for both Read and Write before Read is called when the handshake has not yet completed.
// See SetDeadline, SetReadDeadline, and SetWriteDeadline.
func (qs *QSocket) Read(b []byte) (int, error) {
	if qs.encConn != nil {
		return qs.encConn.Read(b)
	}
	if qs.tlsConn != nil {
		return qs.tlsConn.Read(b)
	}
	if qs.conn != nil {
		return qs.conn.Read(b)
	}
	return 0, ErrUninitializedSocket
}

// Write writes data to the connection.
//
// As Write calls Handshake, in order to prevent indefinite blocking a deadline must be set for both Read and Write before Write is called when the handshake has not yet completed.
// See SetDeadline, SetReadDeadline, and SetWriteDeadline.
func (qs *QSocket) Write(b []byte) (int, error) {
	if qs.encConn != nil {
		return qs.encConn.Write(b)
	}
	if qs.tlsConn != nil {
		return qs.tlsConn.Write(b)
	}
	if qs.conn != nil {
		return qs.conn.Write(b)
	}
	return 0, ErrUninitializedSocket
}

// Close closes the QSocket connection and underlying TCP/TLS connections.
func (qs *QSocket) Close() {
	if qs.encConn != nil {
		qs.encConn.Close()
	}
	if qs.tlsConn != nil {
		qs.tlsConn.Close()
	}
	if qs.conn != nil {
		qs.conn.Close()
	}
	qs.conn = nil
	qs.tlsConn = nil
	qs.encConn = nil
}

// chanFromConn creates a channel from a Conn object, and sends everything it
//
//	Read()s from the socket to the channel.
func CreateSocketChan(sock *QSocket) chan []byte {
	c := make(chan []byte)

	go func() {
		b := make([]byte, 1024)
		for {
			if sock.IsClosed() {
				c <- nil
				return
			}
			sock.SetReadDeadline(time.Time{})
			n, err := sock.Read(b)
			if n > 0 {
				res := make([]byte, n)
				// Copy the buffer so it doesn't get changed while read by the recipient.
				copy(res, b[:n])
				c <- res
			}
			if err != nil || sock.IsClosed() {
				// if err.Error() != "EOF" {
				// 	logrus.Errorf("%s -read-err-> %s", sock.RemoteAddr(), err)
				// }
				c <- nil
				break
			}
		}
	}()

	return c
}

// BindSockets is used for creating a full duplex channel between `con1` and `con2` sockets,
// effectively binding two sockets.
func BindSockets(con1, con2 *QSocket) error {
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
				err = ErrQSocketSessionEnd
			}
		case b2 := <-chan2:
			if b2 != nil {
				_, err = con1.Write(b2)
			} else {
				err = ErrQSocketSessionEnd
			}
		}
		if err != nil {
			break
		}
	}
	return err
}

// GetDefaultTag creates a device ID tag based
// based on the device operating system and architecture.
func GetDefaultTag() byte {
	tag := byte(0)
	switch runtime.GOARCH {
	case "amd64":
		tag = tag | TAG_ARCH_AMD64
	case "386":
		tag = tag | TAG_ARCH_386
	case "arm64":
		tag = tag | TAG_ARCH_ARM64
	case "arm":
		tag = tag | TAG_ARCH_ARM
	case "mips":
		tag = tag | TAG_ARCH_MIPS
	case "mips64":
		tag = tag | TAG_ARCH_MIPS64
	case "mips64le":
		tag = tag | TAG_ARCH_MIPS64LE
	}

	switch runtime.GOOS {
	case "linux":
		tag = tag | TAG_OS_LINUX
	case "windows":
		tag = tag | TAG_OS_WINDOWS
	case "darwin":
		tag = tag | TAG_OS_DARWIN
	case "android":
		tag = tag | TAG_OS_ANDROID
	case "ios":
		tag = tag | TAG_OS_IOS
	case "freebsd":
		tag = tag | TAG_OS_FREEBSD
	case "openbsd":
		tag = tag | TAG_OS_OPENBSD
	}

	return tag
}
