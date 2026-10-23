package goldie

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name       string
		cfg        config
		setup      func(t *testing.T, dir string) string
		wantFailed bool
		wantStale  bool
	}{
		{
			name:  "creates nested dirs",
			cfg:   config{dirPerms: 0o755, filePerms: 0o644},
			setup: func(t *testing.T, dir string) string { return filepath.Join(dir, "a", "b.golden") },
		},
		{
			name: "clean removes stale golden",
			cfg:  config{dirPerms: 0o755, filePerms: 0o644, clean: true},
			setup: func(t *testing.T, dir string) string {
				writeFile(t, filepath.Join(dir, "stale.golden"), "x")
				return filepath.Join(dir, "fresh.golden")
			},
		},
		{
			name: "no clean keeps stale golden",
			cfg:  config{dirPerms: 0o755, filePerms: 0o644},
			setup: func(t *testing.T, dir string) string {
				writeFile(t, filepath.Join(dir, "stale.golden"), "x")
				return filepath.Join(dir, "fresh.golden")
			},
			wantStale: true,
		},
		{
			name: "mkdir failure",
			cfg:  config{dirPerms: 0o755, filePerms: 0o644},
			setup: func(t *testing.T, dir string) string {
				writeFile(t, filepath.Join(dir, "file"), "i am a file")
				return filepath.Join(dir, "file", "a.golden")
			},
			wantFailed: true,
		},
		{
			name: "write failure",
			cfg:  config{dirPerms: 0o755, filePerms: 0o644},
			setup: func(t *testing.T, dir string) string {
				path := filepath.Join(dir, "a.golden")
				require.NoError(t, os.MkdirAll(path, 0o755))
				return path
			},
			wantFailed: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ft := &fakeT{}
			dir := filepath.Join(t.TempDir(), "testdata")
			path := tc.setup(t, dir)

			run(func() { write(ft, path, []byte("data"), tc.cfg) })

			require.Equal(t, tc.wantFailed, ft.failed)
			if tc.wantFailed {
				return
			}
			got, err := os.ReadFile(path)
			require.NoError(t, err)
			require.Equal(t, "data", string(got))
			stale := filepath.Join(dir, "stale.golden")
			switch {
			case tc.wantStale:
				require.FileExists(t, stale)
			case tc.cfg.clean:
				require.NoFileExists(t, stale)
				require.Contains(t, ft.logs, "removed golden file "+stale)
			}
		})
	}
}

func TestCleanDir(t *testing.T) {
	tests := []struct {
		name       string
		ext        string
		setup      func(t *testing.T, dir string)
		wantFailed bool
		wantGone   []string
		wantKept   []string
	}{
		{name: "missing dir is ignored", ext: ".golden", setup: func(t *testing.T, dir string) {}},
		{
			name: "removes only same-extension files, not recursively",
			ext:  ".golden",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "a.golden"), "x")
				writeFile(t, filepath.Join(dir, "input.json"), "x")
				writeFile(t, filepath.Join(dir, "sub", "b.golden"), "x")
			},
			wantGone: []string{"a.golden"},
			wantKept: []string{"input.json", "sub/b.golden"},
		},
		{
			name: "no extension is never cleaned",
			ext:  "",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "Makefile"), "x")
			},
			wantKept: []string{"Makefile"},
		},
		{
			name: "read dir failure",
			ext:  ".golden",
			setup: func(t *testing.T, dir string) {
				writeFile(t, dir, "i am a file")
			},
			wantFailed: true,
		},
		{
			name: "remove failure",
			ext:  ".golden",
			setup: func(t *testing.T, dir string) {
				skipIfRoot(t)
				writeFile(t, filepath.Join(dir, "a.golden"), "x")
				require.NoError(t, os.Chmod(dir, 0o555))
				t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
			},
			wantFailed: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ft := &fakeT{}
			dir := filepath.Join(t.TempDir(), "testdata")
			tc.setup(t, dir)

			run(func() { cleanDir(ft, dir, tc.ext) })

			require.Equal(t, tc.wantFailed, ft.failed)
			for _, n := range tc.wantGone {
				require.NoFileExists(t, filepath.Join(dir, n))
				require.Contains(t, ft.logs, "removed golden file "+filepath.Join(dir, n))
			}
			for _, n := range tc.wantKept {
				require.FileExists(t, filepath.Join(dir, n))
			}
		})
	}

	t.Run("runs once per dir and extension", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "testdata")
		ft := &fakeT{}
		writeFile(t, filepath.Join(dir, "a.golden"), "x")
		cleanDir(ft, dir, ".golden")
		writeFile(t, filepath.Join(dir, "b.golden"), "x")
		cleanDir(ft, dir, ".golden")

		require.NoFileExists(t, filepath.Join(dir, "a.golden"))
		require.FileExists(t, filepath.Join(dir, "b.golden"))
	})
}
