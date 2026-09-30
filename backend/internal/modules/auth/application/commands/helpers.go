package commands

import "encoding/base64"

func encodeOpaque(value []byte) string {
	return base64.RawURLEncoding.EncodeToString(value)
}
