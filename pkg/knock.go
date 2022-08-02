package qsnetcat

import (
	"crypto/md5"
	"errors"

	"github.com/qsocket/qs-netcat/utils"

	"github.com/sirupsen/logrus"
)

const (
	KNOCK_CHECKSUM_BASE      = 0xEE
	KNOCK_HEADER_B1     uint = 0xC0
	KNOCK_HEADER_B2     uint = 0xDE
	KNOCK_SUCCESS       uint = 0xE0
	KNOCK_FAIL          uint = 0xE1
	KNOCK_BUSY          uint = 0xE2
)

var (
	ErrInvalidKnockResponse = errors.New("Invalid knock response!")
	ErrKnockSendFailed      = errors.New("Knock sequence send failed!")
	ErrConnRefused          = errors.New("Connection refused (no server listening with given secret)")
	ErrSocketBusy           = errors.New("Socket busy!")
)

func SendKnockSequence(conn *QuantumSocket, secret string, tag uint8) error {
	knock := []byte{byte(KNOCK_HEADER_B1), byte(KNOCK_HEADER_B2)}
	uid := md5.Sum([]byte(secret))
	checksum := utils.CaclChecksum(uid[:], KNOCK_CHECKSUM_BASE)
	knock = append(knock, byte(checksum))
	knock = append(knock, uid[:]...)
	knock = append(knock, byte(tag))

	n, err := conn.Write(knock)
	if err != nil {
		return err
	}

	if n != 20 {
		return ErrKnockSendFailed
	}

	resp := make([]byte, 1)
	_, err = conn.Read(resp)
	if err != nil {
		return err
	}

	switch resp[0] {
	case byte(KNOCK_SUCCESS):
		return nil
	case byte(KNOCK_BUSY):
		return ErrSocketBusy
	case byte(KNOCK_FAIL):
		return ErrConnRefused
	default:
		logrus.Debugf("Received response: %x", resp[0])
		return ErrInvalidKnockResponse
	}
}
