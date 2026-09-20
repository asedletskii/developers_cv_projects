package valueobjects

import (
	"errors"
	"fmt"
)

const HashSize = 32

var ErrInvalidHash = errors.New("invalid hash")

type SessionTokenHash struct {
	value [HashSize]byte
}

func NewSessionTokenHash(value []byte) (SessionTokenHash, error) {
	if len(value) != HashSize {
		return SessionTokenHash{}, fmt.Errorf("%w: session token hash must be %d bytes", ErrInvalidHash, HashSize)
	}
	var hash [HashSize]byte
	copy(hash[:], value)
	return SessionTokenHash{value: hash}, nil
}

func (h SessionTokenHash) Bytes() []byte {
	value := make([]byte, HashSize)
	copy(value, h.value[:])
	return value
}

func (h SessionTokenHash) IsZero() bool {
	return h.value == [HashSize]byte{}
}
