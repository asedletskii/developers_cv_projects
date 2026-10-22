package contact

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Key
		wantErr bool
	}{
		{name: "phone", input: "phone", want: Phone},
		{name: "telegram", input: "telegram", want: Telegram},
		{name: "email", input: "email", want: Email},
		{name: "unknown", input: "fax", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			require.True(t, got.IsValid())
		})
	}
}
