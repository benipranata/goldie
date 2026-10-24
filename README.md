# goldie

[![Go Reference](https://pkg.go.dev/badge/github.com/benipranata/goldie/pkg/goldie.svg)](https://pkg.go.dev/github.com/benipranata/goldie/pkg/goldie)

**Golden file testing for Go, built on [testify/require](https://github.com/stretchr/testify).**

You write one assertion. goldie compares the output with a file in `testdata/`, shows a real diff when it changes, and rewrites the file with one command.

- **Tiny API:** `goldie.New(t).Equal(name, actual)`. That's it.
- **Built on testify:** failures are `require.Equal` failures, with the diff you already know.
- **Always up to date:** `-update` rewrites your golden files, and `-clean` removes the stale ones.
- **Quiet diffs:** regenerating unchanged output touches no files, so `git diff` stays empty.

```go
func TestRender(t *testing.T) {
	goldie.New(t).Equal("render", []byte(Render())) // testdata/render.golden
}
```

## Contents

- [Install](#install)
- [Quickstart](#quickstart)
- [Examples](#examples)
- [Flags](#flags)
- [Options](#options)
- [Errors](#errors)
- [Limitations](#limitations)
- [Development](#development)
- [License](#license)

## Install

Requires Go 1.22 or later.

```sh
go get github.com/benipranata/goldie
```

## Quickstart

1. Add the assertion to a test:

   ```go
   import "github.com/benipranata/goldie/pkg/goldie"

   func TestRender(t *testing.T) {
   	goldie.New(t).Equal("render", []byte(Render()))
   }
   ```

2. Optionally, add a `go:generate` line so regenerating is one command:

   ```go
   //go:generate go test -v -update -clean .
   ```

3. Create the golden files, then run your tests as usual:

   ```sh
   go generate ./...   # writes testdata/render.golden
   go test ./...       # compares against it
   ```

4. Change the output, and the test shows exactly what moved:

   ```
   --- Expected
   +++ Actual
   @@ -1,2 +1,2 @@
   -Hello, World
   +Hello, Gopher
   ```

5. Happy with the new output? Run `go generate ./...` again and commit the updated file.

## Examples

Subtests, one golden file each (`testdata/TestRender/<subtest>.golden`):

```go
func TestRender(t *testing.T) {
	g := goldie.New(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g.Equal(t.Name(), render(tc.input))
		})
	}
}
```

JSON, in the format you choose:

```go
b, _ := json.MarshalIndent(got, "", "  ")
goldie.New(t).Equal("state", b)
```

HTTP responses:

```go
rec := httptest.NewRecorder()
handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/users", nil))
goldie.New(t).Equal("users", rec.Body.Bytes())
```

## Flags

| Flag | What it does |
|---|---|
| `-update` | Writes the actual output to the golden files instead of comparing. Files whose content is unchanged are left alone. |
| `-clean` | With `-update`, first removes every file with the golden suffix under the fixture directory, once per run. Golden files for deleted tests don't pile up. |

```sh
go test -v -update ./...          # create or refresh
go test -v -update -clean ./...   # refresh and remove the stale ones
```

Add `-v` to see which files were written or removed.

Run `-clean` on the whole package. A test that is skipped, or excluded by a build tag, won't recreate the golden file that `-clean` removed.

## Options

Pass options to `goldie.New`:

```go
g := goldie.New(t,
	goldie.WithFixtureDir("fixtures"),
	goldie.WithNameSuffix(".txt"),
)
```

| Option | Default | What it sets |
|---|---|---|
| `WithFixtureDir(dir)` | `testdata` | Directory holding the golden files. |
| `WithNameSuffix(suffix)` | `.golden` | Suffix appended to each name. An empty suffix is never cleaned. |
| `WithDirPerms(perm)` | `0755` | Permissions of directories goldie creates. |
| `WithFilePerms(perm)` | `0644` | Permissions of golden files goldie writes. |

## Errors

A missing golden file and a changed output both fail the test immediately, and both name the file and the fix:

```
read golden file testdata/render.golden (run with -update to create it)
golden file testdata/render.golden mismatch (run with -update to accept)
```

## Limitations

- Comparison is byte for byte. Normalize timestamps, random IDs and map ordering before calling `Equal`.
- `Equal` takes `[]byte`. Convert strings and encode structs yourself, which also lets you pick the format.
- The golden file name is joined onto the fixture directory as given, so `t.Name()` produces nested directories for subtests.

## Development

```sh
go test ./...
```

| Path | Responsibility |
|---|---|
| `pkg/goldie/goldie.go` | `New`, `Equal`, and update and clean logic. |
| `pkg/goldie/options.go` | `With*` options. |
| `pkg/goldie/flags.go` | The `-update` and `-clean` flags. |
| `pkg/goldie/testdata/` | Golden files used by the integration test. |

## License

[MIT](LICENSE)
