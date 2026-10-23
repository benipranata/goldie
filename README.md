# goldie

Golden file testing built on [testify/require](https://github.com/stretchr/testify), inspired by [sebdah/goldie](https://github.com/sebdah/goldie).

```sh
go get github.com/benipranata/goldie
```

```go
import "github.com/benipranata/goldie/pkg/goldie"

//go:generate go test -v -update -clean .

func TestRender(t *testing.T) {
	goldie.Equal(t, "testdata/render.golden", Render())
}
```

`goldie.Equal(t, expPath, actual, opts...)` requires `actual` (a `string` or `[]byte`) to equal the contents of the golden file at `expPath`. A mismatch fails through `require.Equal`, so you get testify's diff.

```sh
go test ./...                 # compare
go generate ./...             # regenerate via the directive above
go test -v -update -clean .   # the same, by hand
```

Use `-v`: `go test` hides the logs of passing tests, and update mode reports each golden file it writes or removes through `t.Logf`.

## Flags

- `-update` writes or overwrites golden files instead of comparing.
- `-clean` only takes effect with `-update`. Before the first write in a directory, it removes the files in that directory with the same extension as the golden file, once per package run. Subdirectories and files with other extensions are left alone. A golden file without an extension is never cleaned.

## Caveats of `-clean`

- Run it on the whole package, never with `-run`: golden files of tests that don't run are deleted and not recreated.
- The same goes for skipped tests and tests excluded by build tags or OS. Their golden files are removed too.
- Cleaning happens at the first write, before later tests have run. If a test then fails, regenerated files are incomplete. Commit your golden files so `git checkout` can restore them.
- A directory is only cleaned when a test writes into it, so a directory left behind by deleted tests is never cleaned.
- `-update` and `-clean` are registered on import, so the package can't be combined with another library that registers the same flags.

## Options

`WithDirPerms` (default `0755`) and `WithFilePerms` (default `0644`).
