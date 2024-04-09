package cestr

import (
	"testing"
)

func TestRandStringBytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				n: 11,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandStringBytes(tt.args.n); got != tt.want {
				t.Errorf("RandStringBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
