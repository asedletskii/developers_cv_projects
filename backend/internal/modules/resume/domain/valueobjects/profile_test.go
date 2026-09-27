package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewContactsTextRejectsValuesThatNormalizeToEmpty(t *testing.T) {
	empty := "   "

	_, err := NewContactsText(ContactsTextParams{Title: &empty})

	require.ErrorIs(t, err, ErrInvalidValue)
}

func TestNewContactsTextNormalizesAndClonesValues(t *testing.T) {
	title := "  Contacts  "
	subtitle := "  Reach me  "

	value, err := NewContactsText(ContactsTextParams{Title: &title, Subtitle: &subtitle})
	require.NoError(t, err)
	require.Equal(t, "Contacts", *value.Title())
	require.Equal(t, "Reach me", *value.Subtitle())

	clone := value.Clone()
	require.Equal(t, "Contacts", *clone.Title())
	require.Equal(t, "Reach me", *clone.Subtitle())
}
