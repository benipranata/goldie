package goldie

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeT struct {
	failed bool
	errs   []string
	logs   []string
}

func (f *fakeT) Errorf(format string, args ...any) {
	f.errs = append(f.errs, fmt.Sprintf(format, args...))
}
func (f *fakeT) FailNow() { f.failed = true; runtime.Goexit() }
func (f *fakeT) Helper()  {}
func (f *fakeT) Logf(format string, args ...any) {
	f.logs = append(f.logs, fmt.Sprintf(format, args...))
}

// run executes fn on its own goroutine so FailNow's Goexit doesn't kill the test.
func run(fn func()) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	<-done
}

func setFlags(t *testing.T, update, clean bool) {
	t.Helper()
	oldU, oldC := *updateFlag, *cleanFlag
	t.Cleanup(func() { *updateFlag, *cleanFlag = oldU, oldC })
	*updateFlag, *cleanFlag = update, clean
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func skipIfRoot(t *testing.T) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("permission errors can't be provoked as root")
	}
}

func ptr(s string) *string { return &s }
