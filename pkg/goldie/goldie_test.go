package goldie

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	type namedString string

	tests := []struct {
		name       string
		golden     *string
		update     bool
		call       func(t TB, path string)
		wantFailed bool
		wantFile   string
		wantMode   os.FileMode
	}{
		{
			name:     "string match",
			golden:   ptr("hello\n"),
			call:     func(t TB, p string) { Equal(t, p, "hello\n") },
			wantFile: "hello\n",
		},
		{
			name:     "bytes match",
			golden:   ptr("hello\n"),
			call:     func(t TB, p string) { Equal(t, p, []byte("hello\n")) },
			wantFile: "hello\n",
		},
		{
			name:     "named string match",
			golden:   ptr("hello\n"),
			call:     func(t TB, p string) { Equal(t, p, namedString("hello\n")) },
			wantFile: "hello\n",
		},
		{
			name:       "mismatch fails",
			golden:     ptr("hello\n"),
			call:       func(t TB, p string) { Equal(t, p, "bye\n") },
			wantFailed: true,
			wantFile:   "hello\n",
		},
		{
			name:       "missing golden fails",
			call:       func(t TB, p string) { Equal(t, p, "hello\n") },
			wantFailed: true,
		},
		{
			name:     "update creates",
			update:   true,
			call:     func(t TB, p string) { Equal(t, p, "new\n") },
			wantFile: "new\n",
		},
		{
			name:     "update overwrites",
			golden:   ptr("old\n"),
			update:   true,
			call:     func(t TB, p string) { Equal(t, p, []byte("new\n")) },
			wantFile: "new\n",
		},
		{
			name:     "options applied",
			update:   true,
			call:     func(t TB, p string) { Equal(t, p, "new\n", WithFilePerms(0o600)) },
			wantFile: "new\n",
			wantMode: 0o600,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setFlags(t, tc.update, false)
			ft := &fakeT{}
			path := filepath.Join(t.TempDir(), "testdata", "case.golden")
			if tc.golden != nil {
				writeFile(t, path, *tc.golden)
			}

			run(func() { tc.call(ft, path) })

			require.Equal(t, tc.wantFailed, ft.failed)
			if tc.update {
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
