package xgo

import (
	"fmt"
	"log"
	"runtime"

	"go.uber.org/zap"
)

// Go goroutine
func Go(fn func()) {
	go try2(fn, nil)
}

func try2(fn func(), cleaner func()) (ret error) {
	if cleaner != nil {
		defer cleaner()
	}
	defer func() {
		if err := recover(); err != nil {
			_, file, line, _ := runtime.Caller(5)
			log.Print("recover", zap.Any("err", err), zap.String("line", fmt.Sprintf("%s:%d", file, line)))
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
