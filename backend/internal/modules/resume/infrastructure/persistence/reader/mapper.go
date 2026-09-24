package reader

import (
	"encoding/json"
	"fmt"
)

func unmarshal(data []byte, target any) error {
	if len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode read model: %w", err)
	}
	return nil
}

func unmarshalOptional(data []byte, target any) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	return unmarshal(data, target)
}

func ensureArray(data []byte) []byte {
	if len(data) == 0 || string(data) == "null" {
		return []byte("[]")
	}
	return data
}
