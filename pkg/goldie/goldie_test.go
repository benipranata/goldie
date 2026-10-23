package goldie

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestGoldie(t *testing.T, ft *fakeT, opts ...Option) *Goldie {
	t.Helper()
	opts = append([]Option{WithFixtureDir(filepath.Join(t.TempDir(), "testdata"))}, opts...)
	g := New(ft, opts...)
	g.update, g.clean = false, false
	return g
}

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		update     bool
		clean      bool
		opts       []Option
		wantDir    string
		wantSuffix string
	}{
		{name: "defaults", wantDir: "testdata", wantSuffix: ".golden"},
		{name: "flags are read", update: true, clean: true, wantDir: "testdata", wantSuffix: ".golden"},
		{name: "options applied", opts: []Option{WithFixtureDir("fx"), WithNameSuffix(".out")}, wantDir: "fx", wantSuffix: ".out"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setFlags(t, tc.update, tc.clean)

			ft := &fakeT{}
			g := New(ft, tc.opts...)

			require.Same(t, ft, g.t)
			require.Equal(t, tc.wantDir, g.fixtureDir)
			require.Equal(t, tc.wantSuffix, g.suffix)
			require.Equal(t, os.FileMode(0o755), g.dirPerms)
			require.Equal(t, os.FileMode(0o644), g.filePerms)
			require.Equal(t, tc.update, g.update)
			require.Equal(t, tc.clean, g.clean)
		})
	}
}

func TestGoldie_Equal(t *testing.T) {
	tests := []struct {
		name       string
		golden     *string
		update     bool
		opts       []Option
		actual     string
		wantFailed bool
		wantFile   string
		skipped    bool
		wantMode   os.FileMode
	}{
		{name: "match", golden: ptr("hello\n"), actual: "hello\n", wantFile: "hello\n"},
		{name: "mismatch fails", golden: ptr("hello\n"), actual: "bye\n", wantFailed: true, wantFile: "hello\n"},
		{name: "missing golden fails", actual: "hello\n", wantFailed: true},
		{name: "update creates", update: true, actual: "new\n", wantFile: "new\n"},
		{name: "update overwrites", golden: ptr("old\n"), update: true, actual: "new\n", wantFile: "new\n"},
		{name: "update with identical content skips write", golden: ptr("same\n"), update: true, actual: "same\n", wantFile: "same\n", skipped: true},
		{name: "options applied", update: true, opts: []Option{WithFilePerms(0o600)}, actual: "new\n", wantFile: "new\n", wantMode: 0o600},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ft := &fakeT{}
			g := newTestGoldie(t, ft, tc.opts...)
			g.update = tc.update
			path := g.path("case")
			if tc.golden != nil {
				writeFile(t, path, *tc.golden)
			}

			run(func() { g.Equal("case", []byte(tc.actual)) })

			require.Equal(t, tc.wantFailed, ft.failed)
			if tc.update && !tc.skipped {
				require.Equal(t, []string{"updated golden file " + path}, ft.logs)
			} else {
				require.Empty(t, ft.logs)
			}
			if tc.wantFile == "" {
				require.NoFileExists(t, path)
				return
			}
			got, err := os.ReadFile(path)
			require.NoError(t, err)
			require.Equal(t, tc.wantFile, string(got))
			if tc.wantMode != 0 {
				info, err := os.Stat(path)
				require.NoError(t, err)
				require.Equal(t, tc.wantMode, info.Mode().Perm())
			}
		})
	}
}

func TestGoldie_path(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		in   string
		want string
	}{
		{name: "default", in: "a", want: filepath.Join("testdata", "a.golden")},
		{name: "nested name", in: "T/sub", want: filepath.Join("testdata", "T", "sub.golden")},
		{name: "custom", opts: []Option{WithFixtureDir("fx"), WithNameSuffix(".txt")}, in: "a", want: filepath.Join("fx", "a.txt")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := New(&fakeT{}, tc.opts...)
			require.Equal(t, tc.want, g.path(tc.in))
		})
	}
}

func TestGoldie_write(t *testing.T) {
	tests := []struct {
		name       string
		clean      bool
		setup      func(t *testing.T, g *Goldie) string
		wantFailed bool
		wantStale  bool
		wantSkip   bool
	}{
		{
			name:  "creates nested dirs",
			setup: func(t *testing.T, g *Goldie) string { return g.path("a/b") },
		},
		{
			name:  "clean removes stale golden",
			clean: true,
			setup: func(t *testing.T, g *Goldie) string {
				writeFile(t, g.path("stale"), "x")
				return g.path("fresh")
			},
		},
		{
			name: "no clean keeps stale golden",
			setup: func(t *testing.T, g *Goldie) string {
				writeFile(t, g.path("stale"), "x")
				return g.path("fresh")
			},
			wantStale: true,
		},
		{
			name: "identical content is not rewritten",
			setup: func(t *testing.T, g *Goldie) string {
				writeFile(t, g.path("a"), "data")
				return g.path("a")
			},
			wantSkip: true,
		},
		{
			name: "mkdir failure",
			setup: func(t *testing.T, g *Goldie) string {
				writeFile(t, g.fixtureDir, "i am a file")
				return g.path("a")
			},
			wantFailed: true,
		},
		{
			name: "write failure",
			setup: func(t *testing.T, g *Goldie) string {
				require.NoError(t, os.MkdirAll(g.path("a"), 0o755))
				return g.path("a")
			},
			wantFailed: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ft := &fakeT{}
			g := newTestGoldie(t, ft)
			g.clean = tc.clean
			path := tc.setup(t, g)

			run(func() { g.write(path, []byte("data")) })

			require.Equal(t, tc.wantFailed, ft.failed)
			if tc.wantFailed {
				return
			}
			got, err := os.ReadFile(path)
			require.NoError(t, err)
			require.Equal(t, "data", string(got))
			if tc.wantSkip {
				require.Empty(t, ft.logs)
			}
			switch {
			case tc.wantStale:
				require.FileExists(t, g.path("stale"))
			case tc.clean:
				require.NoFileExists(t, g.path("stale"))
				require.Contains(t, ft.logs, "removed golden file "+g.path("stale"))
			}
		})
	}
}

func TestGoldie_cleanFixtures(t *testing.T) {
	tests := []struct {
		name       string
		suffix     string
		setup      func(t *testing.T, g *Goldie)
		wantFailed bool
		wantGone   []string
		wantKept   []string
	}{
		{name: "missing dir is ignored", suffix: ".golden", setup: func(t *testing.T, g *Goldie) {}},
		{
			name:   "removes only suffixed files, recursively",
			suffix: ".golden",
			setup: func(t *testing.T, g *Goldie) {
				writeFile(t, g.path("a"), "x")
				writeFile(t, g.path("sub/b"), "x")
				writeFile(t, filepath.Join(g.fixtureDir, "input.json"), "x")
			},
			wantGone: []string{"a.golden", "sub/b.golden"},
			wantKept: []string{"input.json"},
		},
		{
			name:   "empty suffix is never cleaned",
			suffix: "",
			setup: func(t *testing.T, g *Goldie) {
				writeFile(t, filepath.Join(g.fixtureDir, "Makefile"), "x")
			},
			wantKept: []string{"Makefile"},
		},
		{
			name:   "walk error fails",
			suffix: ".golden",
			setup: func(t *testing.T, g *Goldie) {
				skipIfRoot(t)
				writeFile(t, g.path("sub/b"), "x")
				sub := filepath.Dir(g.path("sub/b"))
				require.NoError(t, os.Chmod(sub, 0o000))
				t.Cleanup(func() { _ = os.Chmod(sub, 0o755) })
			},
			wantFailed: true,
		},
		{
			name:   "remove error fails",
			suffix: ".golden",
			setup: func(t *testing.T, g *Goldie) {
				skipIfRoot(t)
				writeFile(t, g.path("a"), "x")
				require.NoError(t, os.Chmod(g.fixtureDir, 0o555))
				t.Cleanup(func() { _ = os.Chmod(g.fixtureDir, 0o755) })
			},
			wantFailed: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ft := &fakeT{}
			g := newTestGoldie(t, ft, WithNameSuffix(tc.suffix))
			tc.setup(t, g)

			run(g.cleanFixtures)

			require.Equal(t, tc.wantFailed, ft.failed)
			for _, n := range tc.wantGone {
				require.NoFileExists(t, filepath.Join(g.fixtureDir, n))
				require.Contains(t, ft.logs, "removed golden file "+filepath.Join(g.fixtureDir, n))
			}
			for _, n := range tc.wantKept {
				require.FileExists(t, filepath.Join(g.fixtureDir, n))
			}
		})
	}

	t.Run("runs once per fixture dir", func(t *testing.T) {
		g := newTestGoldie(t, &fakeT{})
		writeFile(t, g.path("a"), "x")
		g.cleanFixtures()
		writeFile(t, g.path("b"), "x")
		g.cleanFixtures()

		require.NoFileExists(t, g.path("a"))
		require.FileExists(t, g.path("b"))
	})
}
