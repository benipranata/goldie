package goldie

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithDirPerms(t *testing.T) {
	tests := []struct {
		name string
		perm os.FileMode
	}{
		{name: "0700", perm: 0o700},
		{name: "0775", perm: 0o775},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &config{}
			WithDirPerms(tc.perm)(c)
			require.Equal(t, tc.perm, c.dirPerms)
		})
	}
}

func TestWithFilePerms(t *testing.T) {
	tests := []struct {
		name string
		perm os.FileMode
	}{
		{name: "0600", perm: 0o600},
		{name: "0664", perm: 0o664},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &config{}
			WithFilePerms(tc.perm)(c)
			require.Equal(t, tc.perm, c.filePerms)
		})
	}
}
