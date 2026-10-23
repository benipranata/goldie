package goldie

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithFixtureDir(t *testing.T) {
	tests := []struct{ name, dir string }{
		{name: "custom", dir: "fx"},
		{name: "empty", dir: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &Goldie{}
			WithFixtureDir(tc.dir)(g)
			require.Equal(t, tc.dir, g.fixtureDir)
		})
	}
}

func TestWithNameSuffix(t *testing.T) {
	tests := []struct{ name, suffix string }{
		{name: "custom", suffix: ".out"},
		{name: "empty", suffix: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &Goldie{}
			WithNameSuffix(tc.suffix)(g)
			require.Equal(t, tc.suffix, g.suffix)
		})
	}
}

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
			g := &Goldie{}
			WithDirPerms(tc.perm)(g)
			require.Equal(t, tc.perm, g.dirPerms)
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
			g := &Goldie{}
			WithFilePerms(tc.perm)(g)
			require.Equal(t, tc.perm, g.filePerms)
		})
	}
}
