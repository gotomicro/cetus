package xhttp

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/http2"
	"golang.org/x/net/publicsuffix"
)

func newHttpClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   10 * time.Second, // 减少 Dial Timeout 时间
		KeepAlive: 60 * time.Second,
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DialContext = dialer.DialContext
	t.MaxIdleConns = 200
	t.MaxIdleConnsPerHost = 100
	t.ForceAttemptHTTP2 = true
	t.TLSHandshakeTimeout = 10 * time.Second
	t.ExpectContinueTimeout = 1 * time.Second

	// try HTTP/1.1 -> HTTP/2
	_ = http2.ConfigureTransport(t)

	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	return &http.Client{
		Transport: t,
		Jar:       cookieJar,
	}
}

func NewRestyClient() *resty.Client {
	return resty.NewWithClient(newHttpClient())
}

// ClientOptions ...
type ClientOptions struct {
	RetryAttempts    int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
}

func NewRestyClientWithOptions(options ClientOptions) *resty.Client {
	client := NewRestyClient()
	if options.RetryAttempts > 0 {
		client.SetRetryCount(options.RetryAttempts)
		if options.RetryWaitTime != 0 {
			client.SetRetryWaitTime(options.RetryWaitTime)
		}
		if options.RetryMaxWaitTime != 0 {
			client.SetRetryMaxWaitTime(options.RetryMaxWaitTime)
		}
	}
	return client
}
