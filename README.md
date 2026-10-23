# goldie

[![Go Reference](https://pkg.go.dev/badge/github.com/benipranata/goldie/pkg/goldie.svg)](https://pkg.go.dev/github.com/benipranata/goldie/pkg/goldie)

Golden file testing for Go, built on [testify/require](https://github.com/stretchr/testify).

One call. Real diffs. One command to regenerate.

```go
func TestRender(t *testing.T) {
	goldie.New(t).Equal("render", []byte(Render())) // testdata/render.golden
}
```

## Why goldie

- **Tiny API:** `goldie.New(t).Equal(name, actual)`. That's it.
- **Built on testify:** failures are `require.Equal` failures, with the diff you already know.
- **Always up to date:** one command rewrites your golden files, and `-clean` removes the stale ones.
- **Quiet diffs:** only files that changed are rewritten, so your commits stay clean.
- **Lightweight:** testify is the only dependency.

## Install

```sh
go get github.com/benipranata/goldie
```

## Quick start

```go
import "github.com/benipranata/goldie/pkg/goldie"

//go:generate go test -v -update -clean .

func TestRender(t *testing.T) {
	goldie.New(t).Equal("render", []byte(Render()))
}
```

Create the golden files, then run your tests as usual:

```sh
go generate ./...   # writes testdata/render.golden
go test ./...       # compares against it
```

Change the output, and the test shows exactly what moved:

```
--- Expected
+++ Actual
@@ -1,2 +1,2 @@
-Hello, World
+Hello, Gopher
```

Happy with the new output? Run `go generate ./...` again and commit the updated file.

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
| `-update` | Writes the actual output to the golden files. |
| `-clean` | With `-update`, removes old golden files first, so ones for deleted tests don't pile up. |

```sh
go test -v -update ./...          # create or refresh
go test -v -update -clean ./...   # refresh and remove the stale ones
```

## Options

```go
g := goldie.New(t,
	goldie.WithFixtureDir("fixtures"),  // default: testdata
	goldie.WithNameSuffix(".txt"),      // default: .golden
	goldie.WithDirPerms(0o750),         // default: 0755
	goldie.WithFilePerms(0o600),        // default: 0644
)
```

## License

[MIT](LICENSE)
