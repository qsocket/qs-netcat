package qsutils

import (
	"crypto/tls"
	"errors"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	QsockReadWriteDeadline = 5 // Minutes
)

var (
	ErrUninitializedSocket = errors.New("quantum socket not initiated")
	ErrUnexpectedSocket    = errors.New("unexpected socket type")
)

type QuantumSocket struct {
	conn    net.Conn
	tlsConn *tls.Conn
}

func NewSocket(conn any) (*QuantumSocket, error) {
	qs := QuantumSocket{}
	switch conn.(type) {
	case net.Conn:
		qs.conn = conn.(net.Conn)
		qs.tlsConn = nil
	case *tls.Conn:
		qs.tlsConn = conn.(*tls.Conn)
		qs.conn = nil
	default:
		return nil, ErrUnexpectedSocket
	}
	return &qs, nil
}

func (qs *QuantumSocket) IsClosed() bool {
	return qs.conn == nil && qs.tlsConn == nil
}

func (qs *QuantumSocket) SetReadDeadline(t time.Time) error {
	if qs.tlsConn != nil {
		return qs.tlsConn.SetReadDeadline(t)
	}

	if qs.conn != nil {
		return qs.conn.SetReadDeadline(t)
	}

	return nil
}

func (qs *QuantumSocket) SetWriteDeadline(t time.Time) error {
	if qs.tlsConn != nil {
		return qs.tlsConn.SetWriteDeadline(t)
	}

	if qs.conn != nil {
		return qs.conn.SetWriteDeadline(t)
	}

	return nil
}

func (qs *QuantumSocket) RemoteAddr() net.Addr {
	if qs.tlsConn != nil {
		return qs.tlsConn.RemoteAddr()
	}

	if qs.conn != nil {
		return qs.conn.RemoteAddr()
	}

	return nil
}

func (qs *QuantumSocket) LocalAddr() net.Addr {
	if qs.tlsConn != nil {
		return qs.tlsConn.RemoteAddr()
	}

	if qs.conn != nil {
		return qs.conn.RemoteAddr()
	}

	return nil
}

func (qs *QuantumSocket) TLS() bool {
	return qs.tlsConn != nil
}

func (qs *QuantumSocket) Read(b []byte) (int, error) {
	if qs.tlsConn != nil {
		return qs.tlsConn.Read(b)
	}

	if qs.conn != nil {
		return qs.conn.Read(b)
	}
	return 0, ErrUninitializedSocket
}

func (qs *QuantumSocket) Write(b []byte) (int, error) {
	if qs.tlsConn != nil {
		return qs.tlsConn.Write(b)
	}

	if qs.conn != nil {
		return qs.conn.Write(b)
	}
	return 0, ErrUninitializedSocket
}

func (qs *QuantumSocket) Close() {
	if qs.tlsConn != nil {
		qs.tlsConn.Close()
	}
	if qs.conn != nil {
		qs.conn.Close()
	}
	qs.conn = nil
	qs.tlsConn = nil
}

// chanFromConn creates a channel from a Conn object, and sends everything it
//  Read()s from the socket to the channel.
func CreateSocketChan(sock *QuantumSocket) chan []byte {
	c := make(chan []byte)

	go func() {
		b := make([]byte, 1024)
		for {
			if sock.IsClosed() {
				c <- nil
				return
			}
			sock.SetReadDeadline(time.Now().Add(QsockReadWriteDeadline * time.Minute))
			n, err := sock.Read(b)
			if n > 0 {
				res := make([]byte, n)
				// Copy the buffer so it doesn't get changed while read by the recipient.
				copy(res, b[:n])
				c <- res
			}
			if err != nil || sock.IsClosed() {
				if err.Error() != "EOF" {
					logrus.Errorf("%s -read-err-> %s", sock.RemoteAddr(), err)
				}
				c <- nil
				break
			}
		}
	}()

	return c
}
