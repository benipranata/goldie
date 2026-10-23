package goldie

import "os"

type config struct {
	dirPerms  os.FileMode
	filePerms os.FileMode
	update    bool
	clean     bool
}

func newConfig(opts []Option) config {
	c := config{
		dirPerms:  0o755,
		filePerms: 0o644,
		update:    *updateFlag,
		clean:     *cleanFlag,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}
