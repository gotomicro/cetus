package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gotomicro/cetus/xcheck"
)

func main() {
	a, _ := xcheck.New(
		xcheck.WithMemOpts(0, 0, 0, 0),
	)
	_ = a.EnableMem().Start()
	go func() {
		for i := 0; i < 100; i++ {
			// memoryLeaking()
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Second * 10)
	// reload the config
	_ = a.Apply(xcheck.WithMemOpts(1024, 10, 20, time.Hour))

	time.Sleep(time.Hour)
}

// memoryLeaking 16MB
func memoryLeaking() {
	type T struct {
		v [1 << 20]int
		t *T
	}

	var finalizer = func(t *T) {
		fmt.Println("finalizer called")
	}

	var x, y T

	// 此SetFinalizer函数调用将使x逃逸到堆上。
	runtime.SetFinalizer(&x, finalizer)

	// 下面这行将形成一个包含x和y的循环引用值组。
	// 这有可能造成x和y不可回收。
	x.t, y.t = &y, &x // y也逃逸到了堆上。
}
