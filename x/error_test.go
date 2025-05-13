package x

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE(t *testing.T) {
	type args struct {
		msg    string
		errors []error
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				msg:    "test",
				errors: []error{errors.New("111"), errors.New("222")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := E(tt.args.msg, tt.args.errors...)
			fmt.Println(res)
		})
	}
}
