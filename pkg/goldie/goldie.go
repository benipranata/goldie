// Package goldie compares test output against golden files using testify's
// require package.
//
// Typical use:
//
//	func TestRender(t *testing.T) {
//		goldie.Equal(t, "testdata/render.golden", Render())
//	}
//
// Regenerate golden files with `go test -update <package>`. Add `-clean` to
// first remove existing golden files, so stale ones don't linger:
//
//	//go:generate go test -v -update -clean .
//
// -clean removes, once per package run, the files in the directory of a golden
// file that share its extension (not subdirectories). Run it on the whole
// package: tests that are skipped, or excluded by build tags, won't recreate
// the golden files it removed. Update mode logs what it writes and removes,
// visible with -v.
package goldie

import (
	"os"

	"github.com/stretchr/testify/require"
)

// TB is the subset of testing.TB that Equal needs.
type TB interface {
	require.TestingT
	Helper()
	Logf(format string, args ...any)
}

// Equal requires actual to equal the contents of the golden file at expPath.
// With -update it writes the golden file instead.
func Equal[T ~string | ~[]byte](t TB, expPath string, actual T, opts ...Option) {
	t.Helper()
	cfg := newConfig(opts)

	if cfg.update {
		write(t, expPath, []byte(actual), cfg)
		return
	}

	expected, err := os.ReadFile(expPath)
	require.NoError(t, err, "read golden file %s (run with -update to create it)", expPath)
	require.Equal(t, string(expected), string(actual), "golden file %s mismatch (run with -update to accept)", expPath)
}
