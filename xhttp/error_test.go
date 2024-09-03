package xhttp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ErrFromCode(t *testing.T) {
	t.Run("should get error ok", func(t *testing.T) {
		cases := map[int]error{
			http.StatusForbidden:       ErrForbidden,
			http.StatusNotFound:        ErrNotFound,
			http.StatusBadRequest:      ErrBadRequest,
			http.StatusUnauthorized:    ErrUnauthorized,
			http.StatusTooManyRequests: ErrInternalServerError,
		}

		for code, err := range cases {
			gotErr := ErrFromCode(code)
			assert.Equal(t, gotErr, err)
		}
	})
}
