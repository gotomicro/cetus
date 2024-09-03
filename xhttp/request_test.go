package xhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func TestGet(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/path/to/timeout", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 50)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	mux.HandleFunc("/path/to/ok", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 50)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	type args struct {
		client  *resty.Client
		url     string
		options *GetOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "should get context.canceled error",
			args: args{
				client: NewRestyClient(),
				url:    server.URL + "/path/to/timeout",
				options: &GetOptions{
					Context:     timeoutCtx,
					QueryParams: nil,
					Headers:     nil,
				},
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "should not get error",
			args: args{
				client: NewRestyClient(),
				url:    server.URL + "/path/to/ok",
				options: &GetOptions{
					Context:     context.Background(),
					QueryParams: nil,
					Headers:     nil,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				_, err := Get(tt.args.client, tt.args.url, tt.args.options)
				if tt.wantErr == nil && err != nil {
					t.Errorf("Get() error should be nil, got %v", err)
					return
				}

				if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			},
		)
	}
}

func TestPost(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/path/to/timeout", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 50)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	mux.HandleFunc("/path/to/ok", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 50)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	type args struct {
		client  *resty.Client
		options *PostOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "should get context.canceled error",
			args: args{
				client: NewRestyClient(),
				options: &PostOptions{
					Context: timeoutCtx,
					Url:     server.URL + "/path/to/timeout",
					Headers: nil,
				},
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "should not get error",
			args: args{
				client: NewRestyClient(),
				options: &PostOptions{
					Context: context.Background(),
					Url:     server.URL + "/path/to/ok",
					Headers: nil,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				_, err := Post(tt.args.client, tt.args.options)
				if tt.wantErr == nil && err != nil {
					t.Errorf("Post() error should be nil, got %v", err)
					return
				}

				if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
					t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			},
		)
	}
}
