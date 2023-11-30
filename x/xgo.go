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

// 通过入参 from to 根据 size 对函数调用进行分页查询，闭区间
// 例如：from=0 to=100 size=10
// 则会调用 10 次函数，每次入参为 from=0 to=10, from=10 to=20, from=20 to=30, ..., from=90 to=100
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
