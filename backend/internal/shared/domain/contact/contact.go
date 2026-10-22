package contact

import "fmt"

type Key string

const (
	Phone    Key = "phone"
	Telegram Key = "telegram"
	Email    Key = "email"
)

func (k Key) IsValid() bool {
	switch k {
	case Phone, Telegram, Email:
		return true
	default:
		return false
	}
}

func Parse(value string) (Key, error) {
	key := Key(value)
	if !key.IsValid() {
		return "", fmt.Errorf("invalid contact key %q", value)
	}
	return key, nil
}
