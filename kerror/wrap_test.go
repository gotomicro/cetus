package kerror

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrap(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				err: errors.New("test-1 error"),
			},
			wantErr: "test-1 error\ntest-2 error\ntest-3 error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Wrap("", tt.args.err)
			err = Wrap(err.Error(), errors.New("test-2 error"))
			err = Wrap(err.Error(), errors.New("test-3 error"))
			fmt.Printf("%+v\n", err)
			if (err != nil) && err.Error() != tt.wantErr {
				t.Errorf("Wrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
