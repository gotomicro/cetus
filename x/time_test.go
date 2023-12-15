package x

import (
	"fmt"
	"testing"
)

func TestAsyncCall(t *testing.T) {
	type args struct {
		f func()
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			args: args{
				f: func() {
					fmt.Println("hello world!!!")
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AsyncCall(tt.args.f)
		})
	}
}
