// Package goldie compares test output against golden files using testify's
// require package.
//
// Typical use:
//
//	func TestRender(t *testing.T) {
//		goldie.New(t).Equal("render", []byte(Render())) // testdata/render.golden
//	}
//
// Regenerate golden files with `go test -update <package>`. Add `-clean` to
// first remove existing golden files, so stale ones don't linger:
//
//	//go:generate go test -v -update -clean .
//
// -clean removes, once per package run, every file with the golden suffix under
// the fixture directory. Run it on the whole package: tests that are skipped,
// or excluded by build tags, won't recreate the golden files it removed.
// Update mode logs what it writes and removes, visible with -v.
package goldie

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/stretchr/testify/require"
)

type cleanKey struct{ dir, suffix string }

var cleaned sync.Map // cleanKey -> *sync.Once

// TB is the subset of testing.TB that Goldie needs.
type TB interface {
	require.TestingT
	Helper()
	Logf(format string, args ...any)
}

// Goldie asserts values against golden files.
type Goldie struct {
	t          TB
	fixtureDir string
	suffix     string
	dirPerms   os.FileMode
	filePerms  os.FileMode
	update     bool
	clean      bool
}

// New returns a Goldie that reports failures to t.
func New(t TB, opts ...Option) *Goldie {
	g := &Goldie{
		t:          t,
		fixtureDir: "testdata",
		suffix:     ".golden",
		dirPerms:   0o755,
		filePerms:  0o644,
		update:     *updateFlag,
		clean:      *cleanFlag,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Equal requires actual to equal the golden file called name. With -update it
// writes the golden file instead.
func (g *Goldie) Equal(name string, actual []byte) {
	g.t.Helper()
	path := g.path(name)

	if g.update {
		g.write(path, actual)
		return
	}

	expected, err := os.ReadFile(path)
	require.NoError(g.t, err, "read golden file %s (run with -update to create it)", path)
	require.Equal(g.t, string(expected), string(actual), "golden file %s mismatch (run with -update to accept)", path)
}

func (g *Goldie) path(name string) string {
	return filepath.Join(g.fixtureDir, name+g.suffix)
}

func (g *Goldie) write(path string, data []byte) {
	g.t.Helper()
	if g.clean {
		g.cleanFixtures()
	}
	if old, err := os.ReadFile(path); err == nil && bytes.Equal(old, data) {
		return
	}
	require.NoError(g.t, os.MkdirAll(filepath.Dir(path), g.dirPerms), "create golden dir")
	require.NoError(g.t, os.WriteFile(path, data, g.filePerms), "write golden file %s", path)
	g.t.Logf("updated golden file %s", path)
}

// cleanFixtures removes the golden files under the fixture dir, once per
// process. An empty suffix is never cleaned, since it would match every file.
func (g *Goldie) cleanFixtures() {
	g.t.Helper()
	if g.suffix == "" {
		return
	}
	once, _ := cleaned.LoadOrStore(cleanKey{g.fixtureDir, g.suffix}, &sync.Once{})
	once.(*sync.Once).Do(func() {
		err := filepath.WalkDir(g.fixtureDir, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(d.Name(), g.suffix) {
				return nil
			}
			if err := os.Remove(p); err != nil {
				return err
			}
			g.t.Logf("removed golden file %s", p)
			return nil
		})
		if !errors.Is(err, fs.ErrNotExist) {
			require.NoError(g.t, err, "clean fixture dir %s", g.fixtureDir)
		}
	})
}
