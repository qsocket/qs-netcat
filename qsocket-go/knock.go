package qsocket

import (
	"crypto/md5"
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

// Some global constants for
// These values can be changed for obfuscating the knock protocol
const (
	// QSRN_GATE is the static gate address for the QSocket network.
	QSRN_GATE = "gate.qsocket.io"
	// QSRN_GATE_TLS_PORT Default TLS port for the QSocket gate.
	QSRN_GATE_TLS_PORT = 443
	// QSRN_GATE_PORT Default TCP port for the QSocket gate.
	QSRN_GATE_PORT = 80
	// CERT_FINGERPRINT is the static TLS certificate fingerprint for QSRN_GATE.
	CERT_FINGERPRINT = "32ADEB12BA582C97E157D10699080C1598ECC3793C09D19020EDF51CDC67C145"

	// KNOCK_CHECKSUM_BASE is the constant base value for calculating knock packet checksums.
	KNOCK_CHECKSUM_BASE = 0xEE
	// KNOCK_HEADER_B1 is the first magic byte of the knock packet.
	KNOCK_HEADER_B1 uint = 0xC0
	// KNOCK_HEADER_B2 is the second magic byte of the knock packet.
	KNOCK_HEADER_B2 uint = 0xDE
	// KNOCK_SUCCESS is the knock sequence response code indicating successful connection.
	KNOCK_SUCCESS uint = 0xE0
	// KNOCK_FAIL is the knock sequence response code indicating failed connection.
	KNOCK_FAIL uint = 0xE1
	// KNOCK_BUSY is the knock sequence response code indicating busy connection.
	KNOCK_BUSY uint = 0xE2
)

var (
	ErrInvalidKnockResponse = errors.New("invalid knock response")
	ErrKnockSendFailed      = errors.New("knock sequence send failed")
	ErrConnRefused          = errors.New("connection refused (no server listening with given secret)")
	ErrSocketBusy           = errors.New("socket busy")
)

// SendKnockSequence sends a knock sequence to the QSRN gate
// with the socket properties.
func (qs *QSocket) SendKnockSequence() error {
	uid := md5.Sum([]byte(qs.Secret))
	if govalidator.IsUUID(qs.Secret) {
		u, err := uuid.Parse(qs.Secret)
		if err != nil {
			return err
		}
		uid = u
	}
	knock, err := NewKnockSequence(uid, qs.tag)
	if err != nil {
		return err
	}
	n, err := qs.Write(knock)
	if err != nil {
		return err
	}
	if n != 20 {
		return ErrKnockSendFailed
	}

	resp := make([]byte, 1)
	_, err = qs.Read(resp)
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
		return ErrInvalidKnockResponse
	}
}

// NewKnockSequence generates a new knock packet with given UUID and tag values.
func NewKnockSequence(uuid [16]byte, tag byte) ([]byte, error) {
	knock := []byte{byte(KNOCK_HEADER_B1), byte(KNOCK_HEADER_B2)}
	checksum := CalcChecksum(uuid[:], KNOCK_CHECKSUM_BASE)
	knock = append(knock, byte(checksum))
	knock = append(knock, uuid[:]...)
	knock = append(knock, tag)
	return knock, nil
}

// CalcChecksum calculates the modulus based checksum of the given data,
// modulus base is given in the base variable.
func CalcChecksum(data []byte, base byte) byte {
	checksum := uint32(0)
	for _, n := range data {
		checksum += uint32((n << 2) % base)
	}
	return byte(checksum % uint32(base))
}
