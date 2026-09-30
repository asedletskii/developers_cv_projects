package valueobjects

import "fmt"

type NonceHash struct {
	value [HashSize]byte
}

func NewNonceHash(value []byte) (NonceHash, error) {
	if len(value) != HashSize {
		return NonceHash{}, fmt.Errorf("nonce hash must be %d bytes", HashSize)
	}
	var hash [HashSize]byte
	copy(hash[:], value)
	return NonceHash{value: hash}, nil
}

func (h NonceHash) Bytes() []byte {
	value := make([]byte, HashSize)
	copy(value, h.value[:])
	return value
}

func (h NonceHash) IsZero() bool {
	return h.value == [HashSize]byte{}
}
