package readers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClickBreakpointCondition(t *testing.T) {
	tests := []struct {
		name       string
		breakpoint string
		want       string
		wantErr    bool
	}{
		{name: "extra small", breakpoint: "xs", want: "c.viewport_w < 600"},
		{
			name:       "medium",
			breakpoint: "md",
			want:       "c.viewport_w >= 600 AND c.viewport_w < 1200",
		},
		{name: "large", breakpoint: "lg", want: "c.viewport_w >= 1200"},
		{name: "unsupported", breakpoint: "xl", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := clickBreakpointCondition(test.breakpoint)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}
