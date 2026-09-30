package valueobjects

import "fmt"

type StateHash struct {
	value [HashSize]byte
}

func NewStateHash(value []byte) (StateHash, error) {
	if len(value) != HashSize {
		return StateHash{}, fmt.Errorf("state hash must be %d bytes", HashSize)
	}
	var hash [HashSize]byte
	copy(hash[:], value)
	return StateHash{value: hash}, nil
}

func (h StateHash) Bytes() []byte {
	value := make([]byte, HashSize)
	copy(value, h.value[:])
	return value
}

func (h StateHash) IsZero() bool {
	return h.value == [HashSize]byte{}
}
