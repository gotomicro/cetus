package ksocketio

import (
	"testing"
)

func BenchmarkDecodeNodeSocketIO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = DecodeNodeSocketIO([]byte(`42["message",{"type":"test-a","data":{"type":"test-b"}}]`), reg)
	}
}

func TestDecodeNodeSocketIO(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name        string
		args        args
		wantTyp     string
		wantContent string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:        "test-1",
			args:        args{in: []byte(`42["message",{"type":"test-a","data":{"type":"test-b"}}]`)},
			wantTyp:     "test-a",
			wantContent: `{"type":"test-a","data":{"type":"test-b"}}`,
			wantErr:     false,
		},
		{
			name:        "test-2",
			args:        args{in: []byte(`0{"type":"test-c"}`)},
			wantTyp:     "test-c",
			wantContent: `{"type":"test-c"}`,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, gotTyp, err := DecodeNodeSocketIO(tt.args.in, reg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeMessageType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTyp != tt.wantTyp {
				t.Errorf("DecodeMessageType() gotTyp = %v, want %v", gotTyp, tt.wantTyp)
			}
			if gotTyp != tt.wantTyp {
				t.Errorf("DecodeMessageType() gotContent = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}
