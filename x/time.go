package x

import (
	"context"
	"fmt"
	"time"
)

func AsyncCall(f func()) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()
	go func(ctx context.Context) {
		// 发送HTTP请求
		f()
		cancel()
	}(ctx)

	select {
	case <-ctx.Done():
		fmt.Println("call successfully!!!")
		return
	case <-time.After(time.Duration(time.Second * 5)):
		fmt.Println("timeout!!!")
		return
	}
}
