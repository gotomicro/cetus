package xmiddleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/cetus/e"
	"golang.org/x/time/rate"

	"github.com/gotomicro/cetus/m/pkg"
)

type (
	// RateLimiterStore is the interface to be implemented by custom stores.
	RateLimiterStore interface {
		// Stores for the rate limiter have to implement the Allow method
		Allow(identifier string) (bool, error)
	}
)

type (
	// RateLimiterConfig defines the configuration for the rate limiter
	RateLimiterConfig struct {
		Skipper    pkg.Skipper
		BeforeFunc pkg.BeforeFunc
		// IdentifierExtractor uses *gin.Context to extract the identifier for a visitor
		IdentifierExtractor Extractor
		// Store defines a store for the rate limiter
		Store RateLimiterStore
		// ErrorHandler provides a handler to be called when IdentifierExtractor returns an error
		ErrorHandler func(context *gin.Context, err error) error
		// DenyHandler provides a handler to be called when RateLimiter denies access
		DenyHandler func(context *gin.Context, identifier string, err error) error
	}
	// Extractor is used to extract data from *gin.Context
	Extractor func(context *gin.Context) (string, error)
)

// errors
var (
	// ErrRateLimitExceeded denotes an error raised when rate limit is exceeded
	ErrRateLimitExceeded = e.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
	// ErrExtractorError denotes an error raised when extractor function is unsuccessful
	ErrExtractorError = e.NewHTTPError(http.StatusForbidden, "error while extracting identifier")
)

// DefaultRateLimiterConfig defines default values for RateLimiterConfig
var DefaultRateLimiterConfig = RateLimiterConfig{
	Skipper: pkg.DefaultSkipper,
	IdentifierExtractor: func(ctx *gin.Context) (string, error) {
		id := ctx.ClientIP()
		return id, nil
	},
	ErrorHandler: func(context *gin.Context, err error) error {
		return &e.HTTPError{
			Code:     ErrExtractorError.Code,
			Message:  ErrExtractorError.Message,
			Internal: err,
		}
	},
	DenyHandler: func(context *gin.Context, identifier string, err error) error {
		return &e.HTTPError{
			Code:     ErrRateLimitExceeded.Code,
			Message:  ErrRateLimitExceeded.Message,
			Internal: err,
		}
	},
}

/*
RateLimiter returns a rate limiting middleware

	limiterStore := middleware.NewRateLimiterMemoryStore(20)
	e.GET("/rate-limited", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
		return
	}, RateLimiter(limiterStore))
*/
func RateLimiter(store RateLimiterStore) gin.HandlerFunc {
	config := DefaultRateLimiterConfig
	config.Store = store

	return RateLimiterWithConfig(config)
}

func RateLimiterWithConfig(config RateLimiterConfig) gin.HandlerFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultRateLimiterConfig.Skipper
	}
	if config.IdentifierExtractor == nil {
		config.IdentifierExtractor = DefaultRateLimiterConfig.IdentifierExtractor
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = DefaultRateLimiterConfig.ErrorHandler
	}
	if config.DenyHandler == nil {
		config.DenyHandler = DefaultRateLimiterConfig.DenyHandler
	}
	if config.Store == nil {
		panic("Store configuration must be provided")
	}
	return func(c *gin.Context) {
		if config.Skipper(c) {
			c.Next()
			return
		}
		if config.BeforeFunc != nil {
			config.BeforeFunc(c)
		}

		identifier, err := config.IdentifierExtractor(c)
		if err != nil {
			c.Error(config.ErrorHandler(c, err))
			return
		}

		if allow, err := config.Store.Allow(identifier); !allow {
			c.Error(config.DenyHandler(c, identifier, err))
			return
		}
		c.Next()
		return
	}
}

type (
	// RateLimiterMemoryStore is the built-in store implementation for RateLimiter
	RateLimiterMemoryStore struct {
		visitors    map[string]*Visitor
		mutex       sync.Mutex
		rate        rate.Limit
		burst       int
		expiresIn   time.Duration
		lastCleanup time.Time
	}
	// Visitor signifies a unique user's limiter details
	Visitor struct {
		*rate.Limiter
		lastSeen time.Time
	}
)

func NewRateLimiterMemoryStore(rate rate.Limit) (store *RateLimiterMemoryStore) {
	return NewRateLimiterMemoryStoreWithConfig(RateLimiterMemoryStoreConfig{
		Rate: rate,
	})
}

func NewRateLimiterMemoryStoreWithConfig(config RateLimiterMemoryStoreConfig) (store *RateLimiterMemoryStore) {
	store = &RateLimiterMemoryStore{}

	store.rate = config.Rate
	store.burst = config.Burst
	store.expiresIn = config.ExpiresIn
	if config.ExpiresIn == 0 {
		store.expiresIn = DefaultRateLimiterMemoryStoreConfig.ExpiresIn
	}
	if config.Burst == 0 {
		store.burst = int(config.Rate)
	}
	store.visitors = make(map[string]*Visitor)
	store.lastCleanup = now()
	return
}

// RateLimiterMemoryStoreConfig represents configuration for RateLimiterMemoryStore
type RateLimiterMemoryStoreConfig struct {
	Rate      rate.Limit    // Rate of requests allowed to pass as req/s
	Burst     int           // Burst additionally allows a number of requests to pass when rate limit is reached
	ExpiresIn time.Duration // ExpiresIn is the duration after that a rate limiter is cleaned up
}

// DefaultRateLimiterMemoryStoreConfig provides default configuration values for RateLimiterMemoryStore
var DefaultRateLimiterMemoryStoreConfig = RateLimiterMemoryStoreConfig{
	ExpiresIn: 3 * time.Minute,
}

// Allow implements RateLimiterStore.Allow
func (store *RateLimiterMemoryStore) Allow(identifier string) (bool, error) {
	store.mutex.Lock()
	limiter, exists := store.visitors[identifier]
	if !exists {
		limiter = new(Visitor)
		limiter.Limiter = rate.NewLimiter(store.rate, store.burst)
		store.visitors[identifier] = limiter
	}
	limiter.lastSeen = now()
	if now().Sub(store.lastCleanup) > store.expiresIn {
		store.cleanupStaleVisitors()
	}
	store.mutex.Unlock()
	return limiter.AllowN(now(), 1), nil
}

/*
cleanupStaleVisitors helps manage the size of the visitors map by removing stale records
of users who haven't visited again after the configured expiry time has elapsed
*/
func (store *RateLimiterMemoryStore) cleanupStaleVisitors() {
	for id, visitor := range store.visitors {
		if now().Sub(visitor.lastSeen) > store.expiresIn {
			delete(store.visitors, id)
		}
	}
	store.lastCleanup = now()
}

/*
actual time method which is mocked in test file
*/
var now = func() time.Time {
	return time.Now()
}
