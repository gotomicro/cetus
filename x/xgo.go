package x

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gotomicro/cetus/l"
)

// Go goroutine
func Go(fn func()) {
	go func() { _ = try(fn, nil) }()
}

func try(fn func(), cleaner func()) (ret error) {
	if cleaner != nil {
		defer cleaner()
	}
	defer func() {
		if err := recover(); err != nil {
			_, file, line, _ := runtime.Caller(5)
			log.Print("recover", l.A("err", err), l.S("line", fmt.Sprintf("%s:%d", file, line)))
			if _, ok := err.(error); ok {
				ret = err.(error)
			} else {
				ret = fmt.Errorf("%+v", err)
			}
		}
	}()
	fn()
	return nil
}
