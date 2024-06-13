package x

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

var ErrAsyncFuncTimeout = errors.New("async func timeout")

type AsyncFunc func(ctx context.Context) error

const CoTimeout = time.Second

// Co Description of function usage
//
//	errCo := x.Co(func(c context.Context) error {
//		return f.xxx(ctx, aa, bb)
//	}, x.CoTimeout)
//	defer func() {
//		if errCo != nil {
//			if errors.Is(errCo, x.ErrAsyncFuncTimeout) {
//				elog.Warn("file_vo_fill_tags", elog.FieldCtxTid(ctx.Request.Context()), l.E(err))
//			} else {
//				elog.Error("file_vo_fill_tags", elog.FieldCtxTid(ctx.Request.Context()), l.E(err))
//			}
//		}
//	}()
func Co(asyncFunc AsyncFunc, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- asyncFunc(ctx)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return errors.New("函数执行超时")
	}
}
