package goldie

import "os"

// Option configures a Goldie.
type Option func(*Goldie)

// WithFixtureDir sets the directory holding golden files. Default: "testdata".
func WithFixtureDir(dir string) Option {
	return func(g *Goldie) { g.fixtureDir = dir }
}

// WithNameSuffix sets the golden file suffix. Default: ".golden".
func WithNameSuffix(suffix string) Option {
	return func(g *Goldie) { g.suffix = suffix }
}

// WithDirPerms sets the permissions of created directories. Default: 0755.
func WithDirPerms(perm os.FileMode) Option {
	return func(g *Goldie) { g.dirPerms = perm }
}

// WithFilePerms sets the permissions of written golden files. Default: 0644.
func WithFilePerms(perm os.FileMode) Option {
	return func(g *Goldie) { g.filePerms = perm }
}
