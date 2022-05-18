package kkafka

import (
	"reflect"
	"testing"

	k "github.com/segmentio/kafka-go"
)

func TestHeadersBatchAdd(t *testing.T) {
	type args struct {
		headers        []k.Header
		needAddHeaders []k.Header
	}
	tests := []struct {
		name string
		args args
		want []k.Header
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			args: args{
				headers: []k.Header{{
					Key:   "A",
					Value: nil,
				}, {
					Key:   "B",
					Value: nil,
				}, {
					Key:   "C",
					Value: nil,
				}},
				needAddHeaders: []k.Header{{
					Key:   "A",
					Value: nil,
				}, {
					Key:   "B",
					Value: nil,
				}, {
					Key:   "D",
					Value: nil,
				}},
			},
			want: []k.Header{{
				Key:   "A",
				Value: nil,
			}, {
				Key:   "B",
				Value: nil,
			}, {
				Key:   "D",
				Value: nil,
			}, {
				Key:   "C",
				Value: nil,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HeadersBatchAdd(tt.args.headers, tt.args.needAddHeaders); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HeadersBatchAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}
