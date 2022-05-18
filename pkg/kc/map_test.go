package kc

import (
	"reflect"
	"testing"
)

func TestKvUnFormat(t *testing.T) {
	type args struct {
		in string
	}
	var tests = []struct {
		name string
		args args
		want []Kv
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			args: args{
				in: "TCP:9001,TCP:9003",
			},
			want: []Kv{
				{
					Key:   "TCP",
					Value: "9001",
				},
				{
					Key:   "TCP",
					Value: "9003",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KvUnFormat(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KvUnFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKv2MapInt(t *testing.T) {
	type args struct {
		in []Kv
	}
	tests := []struct {
		name    string
		args    args
		wantOut map[int]int
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			args: args{
				in: []Kv{
					{
						Key:   2,
						Value: 1,
					},
				},
			},
			wantOut: map[int]int{2: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := Kv2MapInt(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("Kv2MapInt() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
