package goldie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name   string
		update bool
		clean  bool
		opts   []Option
		want   config
	}{
		{name: "defaults", want: config{dirPerms: 0o755, filePerms: 0o644}},
		{name: "flags are read", update: true, clean: true, want: config{dirPerms: 0o755, filePerms: 0o644, update: true, clean: true}},
		{
			name: "options applied",
			opts: []Option{WithDirPerms(0o700), WithFilePerms(0o600)},
			want: config{dirPerms: 0o700, filePerms: 0o600},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setFlags(t, tc.update, tc.clean)
			require.Equal(t, tc.want, newConfig(tc.opts))
		})
	}
}
