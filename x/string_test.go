package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputCheck(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				input: "abc",
			},
			want: true,
		},
		{
			name: "test1",
			args: args{
				input: "abc'~",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, InputCheck(tt.args.input), "InputCheck(%v)", tt.args.input)
		})
	}
}
