package valueobjects

import (
	"fmt"
	"net/mail"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)
	parsed, err := mail.ParseAddress(value)
	if err != nil || parsed.Address != value || !strings.Contains(value, "@") {
		return Email{}, fmt.Errorf("%w: invalid email", ErrInvalidValue)
	}
	return Email{value: value}, nil
}

func (e Email) String() string { return e.value }

type Telegram struct {
	value string
}

func NewTelegram(value string) (Telegram, error) {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "@")
	if value == "" {
		return Telegram{}, fmt.Errorf("%w: telegram is empty", ErrInvalidValue)
	}
	for _, r := range value {
		if !isValidTelegramRune(r) {
			return Telegram{}, fmt.Errorf("%w: invalid telegram username", ErrInvalidValue)
		}
	}
	return Telegram{value: value}, nil
}

func isValidTelegramRune(r rune) bool {
	return r == '_' ||
		(r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}

func (t Telegram) String() string { return t.value }
