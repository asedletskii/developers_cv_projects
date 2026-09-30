package valueobjects

import "fmt"

const HashSize = 32

type SessionTokenHash struct {
	value [HashSize]byte
}

func NewSessionTokenHash(value []byte) (SessionTokenHash, error) {
	if len(value) != HashSize {
		return SessionTokenHash{}, fmt.Errorf("session token hash must be %d bytes", HashSize)
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
