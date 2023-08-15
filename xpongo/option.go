package xpongo

import (
	"io/fs"
)

type config struct {
	path          string
	fs            fs.FS
	debug         bool
	globalContext map[string]interface{}
}

type Option func(opt *config)

func WithPath(path string) Option {
	return func(opt *config) {
		opt.path = path
	}
}

func WithFS(fs fs.FS) Option {
	return func(opt *config) {
		opt.fs = fs
	}
}

func WithDebug(debug bool) Option {
	return func(opt *config) {
		opt.debug = debug
	}
}

func WithGlobalContext(globalContext map[string]interface{}) Option {
	return func(opt *config) {
		opt.globalContext = globalContext
	}
}
