package valueobjects

import "fmt"

type BrowserBindingHash struct {
	value [HashSize]byte
}

func NewBrowserBindingHash(value []byte) (BrowserBindingHash, error) {
	if len(value) != HashSize {
		return BrowserBindingHash{}, fmt.Errorf("browser binding hash must be %d bytes", HashSize)
	}
	var hash [HashSize]byte
	copy(hash[:], value)
	return BrowserBindingHash{value: hash}, nil
}

func (h BrowserBindingHash) Bytes() []byte {
	value := make([]byte, HashSize)
	copy(value, h.value[:])
	return value
}

func (h BrowserBindingHash) IsZero() bool {
	return h.value == [HashSize]byte{}
}
