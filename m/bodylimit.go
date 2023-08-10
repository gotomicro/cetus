package m

import (
	"fmt"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/bytes"

	"github.com/gotomicro/cetus/m/pkg"
)

type (
	// BodyLimitConfig defines the config for BodyLimit middleware.
	BodyLimitConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper pkg.Skipper

		// Maximum allowed size for a request body, it can be specified
		// as `4x` or `4xB`, where x is one of the multiple from K, M, G, T or P.
		Limit string `yaml:"limit"`
		limit int64
	}

	limitedReader struct {
		BodyLimitConfig
		reader  io.ReadCloser
		read    int64
		context *gin.Context
	}
)

var (
	// DefaultBodyLimitConfig is the default BodyLimit middleware config.
	DefaultBodyLimitConfig = BodyLimitConfig{
		Skipper: pkg.DefaultSkipper,
	}
)

// BodyLimit returns a BodyLimit middleware.
//
// BodyLimit middleware sets the maximum allowed size for a request body, if the
// size exceeds the configured limit, it sends "413 - Request Entity Too Large"
// response. The BodyLimit is determined based on both `Content-Length` request
// header and actual content read, which makes it super secure.
// Limit can be specified as `4x` or `4xB`, where x is one of the multiple from K, M,
// G, T or P.
func BodyLimit(limit string) gin.HandlerFunc {
	c := DefaultBodyLimitConfig
	c.Limit = limit
	return BodyLimitWithConfig(c)
}

// BodyLimitWithConfig returns a BodyLimit middleware with config.
// See: `BodyLimit()`.
func BodyLimitWithConfig(config BodyLimitConfig) gin.HandlerFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultBodyLimitConfig.Skipper
	}

	limit, err := bytes.Parse(config.Limit)
	if err != nil {
		panic(fmt.Errorf("echo: invalid body-limit=%s", config.Limit))
	}
	config.limit = limit
	pool := limitedReaderPool(config)

	return func(c *gin.Context) {
		if config.Skipper(c) {
			return next(c)
		}

		req := c.Request()

		// Based on content length
		if req.ContentLength > config.limit {
			return echo.ErrStatusRequestEntityTooLarge
		}

		// Based on content read
		r := pool.Get().(*limitedReader)
		r.Reset(req.Body, c)
		defer pool.Put(r)
		req.Body = r
		c.Next()
		return
	}
}

func (r *limitedReader) Read(b []byte) (n int, err error) {
	n, err = r.reader.Read(b)
	r.read += int64(n)
	if r.read > r.limit {
		return n, echo.ErrStatusRequestEntityTooLarge
	}
	return
}

func (r *limitedReader) Close() error {
	return r.reader.Close()
}

func (r *limitedReader) Reset(reader io.ReadCloser, context *gin.Context) {
	r.reader = reader
	r.context = context
	r.read = 0
}

func limitedReaderPool(c BodyLimitConfig) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return &limitedReader{BodyLimitConfig: c}
		},
	}
}
