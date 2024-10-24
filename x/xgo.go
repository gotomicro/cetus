package x

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gotomicro/cetus/l"
	"github.com/pkg/errors"
)

// Go goroutine
func Go(fn func()) {
	go func() { _ = try(fn, nil) }()
}

func Recover() (err error) {
	if r := recover(); r != nil {
		es := fmt.Sprintf("recover from: %v", r)
		pc, file, line, ok := runtime.Caller(3) // 3 表示向上回溯 3 层
		if ok {
			err = errors.New(fmt.Sprintf("%s, panic occurred in %s[%s:%d]\n", es, runtime.FuncForPC(pc).Name(), file, line))
		} else {
			err = errors.New(es)
		}
		return err
	}
	return nil
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

// Page parses the function call by size from to, with closed intervals.
// For example: from=0 to=100 size=10
// then the function is called 10 times, each time with from=0 to=10, from=10 to=20, from=20 to=30, ... , from=90 to=100
func Page(from, to, size int64, fn func(from, to int64) error) error {
	if from > to {
		return fmt.Errorf("from must less than to")
	}
	if size <= 0 {
		return fmt.Errorf("size must greater than 0")
	}
	for i := from; i <= to; i += size {
		if err := fn(i, i+size-1); err != nil {
			return err
		}
	}
	return nil
}

func SafeFunc(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("Recovered from panic:", r))
		}
	}()
	fn()
	return nil
}
