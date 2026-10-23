package goldie

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/testify/require"
)

type cleanKey struct{ dir, ext string }

var cleaned sync.Map // cleanKey -> *sync.Once

func write(t TB, path string, data []byte, cfg config) {
	t.Helper()
	if cfg.clean {
		cleanDir(t, filepath.Dir(path), filepath.Ext(path))
	}
	require.NoError(t, os.MkdirAll(filepath.Dir(path), cfg.dirPerms), "create golden dir")
	require.NoError(t, os.WriteFile(path, data, cfg.filePerms), "write golden file %s", path)
	t.Logf("updated golden file %s", path)
}

// cleanDir removes the files in dir with extension ext, once per process. A
// golden file without an extension is never cleaned, since ext would match
// every extensionless file.
func cleanDir(t TB, dir, ext string) {
	t.Helper()
	if ext == "" {
		return
	}
	once, _ := cleaned.LoadOrStore(cleanKey{dir, ext}, &sync.Once{})
	once.(*sync.Once).Do(func() {
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			return
		}
		require.NoError(t, err, "clean golden dir %s", dir)

		for _, e := range entries {
			if e.IsDir() || filepath.Ext(e.Name()) != ext {
				continue
			}
			p := filepath.Join(dir, e.Name())
			require.NoError(t, os.Remove(p), "remove golden file %s", p)
			t.Logf("removed golden file %s", p)
		}
	})
}
