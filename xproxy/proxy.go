package xproxy

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/cetus/e"
)

type (
	// ProxyConfig defines the config for Proxy middleware.
	ProxyConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Balancer defines a load balancing technique.
		// Required.
		Balancer ProxyBalancer

		// Rewrite defines URL path rewrite rules. The values captured in asterisk can be
		// retrieved by index e.g. $1, $2 and so on.
		// Examples:
		// "/old":              "/new",
		// "/api/*":            "/$1",
		// "/js/*":             "/public/javascripts/$1",
		// "/users/*/orders/*": "/user/$1/order/$2",
		Rewrite map[string]string

		// RegexRewrite defines rewrite rules using regexp.Rexexp with captures
		// Every capture group in the values can be retrieved by index e.g. $1, $2 and so on.
		// Example:
		// "^/old/[0.9]+/":     "/new",
		// "^/api/.+?/(.*)":    "/v2/$1",
		RegexRewrite map[*regexp.Regexp]string

		// Context key to store selected ProxyTarget into context.
		// Optional. Default value "target".
		ContextKey string

		// To customize the transport to remote.
		// Examples: If custom TLS certificates are required.
		Transport http.RoundTripper

		// ModifyResponse defines function to modify response from ProxyTarget.
		ModifyResponse func(*http.Response) error
	}

	// ProxyTarget defines the upstream target.
	ProxyTarget struct {
		Name string
		URL  *url.URL
		Meta map[string]interface{}
	}

	// ProxyBalancer defines an interface to implement a load balancing technique.
	ProxyBalancer interface {
		AddTarget(*ProxyTarget) bool
		RemoveTarget(string) bool
		Next(*gin.Context) *ProxyTarget
	}

	commonBalancer struct {
		targets []*ProxyTarget
		mutex   sync.RWMutex
	}

	// RandomBalancer implements a random load balancing technique.
	randomBalancer struct {
		*commonBalancer
		random *rand.Rand
	}

	// RoundRobinBalancer implements a round-robin load balancing technique.
	roundRobinBalancer struct {
		*commonBalancer
		i uint32
	}
)

var (
	// DefaultProxyConfig is the default Proxy middleware config.
	DefaultProxyConfig = ProxyConfig{
		Skipper:    DefaultSkipper,
		ContextKey: "target",
	}
)

func proxyRaw(t *ProxyTarget, c *gin.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, _, err := c.Writer.(http.Hijacker).Hijack()
		if err != nil {
			c.Set("_error", fmt.Sprintf("proxy raw, hijack error=%v, url=%s", t.URL, err))
			return
		}
		defer in.Close()

		out, err := net.Dial("tcp", t.URL.Host)
		if err != nil {
			c.Set("_error", e.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("proxy raw, dial error=%v, url=%s", t.URL, err)))
			return
		}
		defer out.Close()

		// Write header
		err = r.Write(out)
		if err != nil {
			c.Set("_error", e.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("proxy raw, request header copy error=%v, url=%s", t.URL, err)))
			return
		}

		errCh := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err = io.Copy(dst, src)
			errCh <- err
		}

		go cp(out, in)
		go cp(in, out)
		err = <-errCh
		if err != nil && err != io.EOF {
			c.Set("_error", fmt.Errorf("proxy raw, copy body error=%v, url=%s", t.URL, err))
		}
	})
}

// NewRandomBalancer returns a random proxy balancer.
func NewRandomBalancer(targets []*ProxyTarget) ProxyBalancer {
	b := &randomBalancer{commonBalancer: new(commonBalancer)}
	b.targets = targets
	return b
}

// NewRoundRobinBalancer returns a round-robin proxy balancer.
func NewRoundRobinBalancer(targets []*ProxyTarget) ProxyBalancer {
	b := &roundRobinBalancer{commonBalancer: new(commonBalancer)}
	b.targets = targets
	return b
}

// AddTarget adds an upstream target to the list.
func (b *commonBalancer) AddTarget(target *ProxyTarget) bool {
	for _, t := range b.targets {
		if t.Name == target.Name {
			return false
		}
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.targets = append(b.targets, target)
	return true
}

// RemoveTarget removes an upstream target from the list.
func (b *commonBalancer) RemoveTarget(name string) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for i, t := range b.targets {
		if t.Name == name {
			b.targets = append(b.targets[:i], b.targets[i+1:]...)
			return true
		}
	}
	return false
}

// Next randomly returns an upstream target.
func (b *randomBalancer) Next(c *gin.Context) *ProxyTarget {
	if b.random == nil {
		b.random = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	}
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.targets[b.random.Intn(len(b.targets))]
}

// Next returns an upstream target using round-robin technique.
func (b *roundRobinBalancer) Next(c *gin.Context) *ProxyTarget {
	b.i = b.i % uint32(len(b.targets))
	t := b.targets[b.i]
	atomic.AddUint32(&b.i, 1)
	return t
}

// Proxy returns a Proxy middleware.
//
// Proxy middleware forwards the request to upstream server using a configured load balancing technique.
func Proxy(balancer ProxyBalancer) gin.HandlerFunc {
	c := DefaultProxyConfig
	c.Balancer = balancer
	return ProxyWithConfig(c)
}

// ProxyWithConfig returns a Proxy middleware with config.
// See: `Proxy()`
func ProxyWithConfig(config ProxyConfig) gin.HandlerFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultProxyConfig.Skipper
	}
	if config.Balancer == nil {
		panic("proxy middleware requires balancer")
	}
	if config.Rewrite != nil {
		if config.RegexRewrite == nil {
			config.RegexRewrite = make(map[*regexp.Regexp]string)
		}
		for k, v := range rewriteRulesRegex(config.Rewrite) {
			config.RegexRewrite[k] = v
		}
	}
	return func(c *gin.Context) {
		if config.Skipper(c) {
			c.Next()
			return
		}

		req := c.Request
		tgt := config.Balancer.Next(c)
		c.Set(config.ContextKey, tgt)

		// Set rewrite path and raw path
		rewritePath(config.RegexRewrite, req)

		// Fix header
		// Basically it's not good practice to unconditionally pass incoming x-real-ip header to upstream.
		// However, for backward compatibility, legacy behavior is preserved unless you configure Echo#IPExtractor.
		if req.Header.Get(HeaderXRealIP) == "" || c.ClientIP() != "" {
			req.Header.Set(HeaderXRealIP, c.ClientIP())
		}
		if req.Header.Get(HeaderXForwardedProto) == "" {
			req.Header.Set(HeaderXForwardedProto, c.Request.URL.Scheme)
		}
		if IsWebSocket(c) && req.Header.Get(HeaderXForwardedFor) == "" { // For HTTP, it is automatically set by Go HTTP reverse proxy.
			req.Header.Set(HeaderXForwardedFor, c.ClientIP())
		}

		// Proxy
		switch {
		case IsWebSocket(c):
			proxyRaw(tgt, c).ServeHTTP(c.Writer, req)
		case req.Header.Get(HeaderAccept) == "text/event-stream":
		default:
			proxyHTTP(tgt, c, config).ServeHTTP(c.Writer, req)
		}
		if err, ok := c.Get("_error"); ok {
			_ = c.Error(err.(error))
			return
		}
		return
	}
}

func IsWebSocket(c *gin.Context) bool {
	upgrade := c.Request.Header.Get(HeaderUpgrade)
	return strings.ToLower(upgrade) == "websocket"
}

func proxyHTTP(tgt *ProxyTarget, c *gin.Context, config ProxyConfig) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(tgt.URL)
	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		desc := tgt.URL.String()
		if tgt.Name != "" {
			desc = fmt.Sprintf("%s(%s)", tgt.Name, tgt.URL.String())
		}
		c.Set("_error", e.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("remote %s unreachable, could not forward: %v", desc, err)))
	}
	proxy.Transport = config.Transport
	proxy.ModifyResponse = config.ModifyResponse
	return proxy
}
