package goldie

import "os"

// Option configures Equal.
type Option func(*config)

// WithDirPerms sets the permissions of created directories. Default: 0755.
func WithDirPerms(perm os.FileMode) Option {
	return func(c *config) { c.dirPerms = perm }
}

// WithFilePerms sets the permissions of written golden files. Default: 0644.
func WithFilePerms(perm os.FileMode) Option {
	return func(c *config) { c.filePerms = perm }
}
