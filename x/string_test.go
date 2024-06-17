package x

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAESEncrypt(t *testing.T) {
	type args struct {
		text string
		key  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name:    "test1",
			args:    args{text: "hello", key: "0011223344556677"},
			want:    "10OwHTs=",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESEncrypt(tt.args.text, tt.args.key)
			if !tt.wantErr(t, err, fmt.Sprintf("AESEncrypt(%v, %v)", tt.args.text, tt.args.key)) {
				return
			}
			assert.Equalf(t, tt.want, got, "AESEncrypt(%v, %v)", tt.args.text, tt.args.key)
		})
	}
}

func TestAESDecrypt(t *testing.T) {
	type args struct {
		text string
		key  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name:    "test1",
			args:    args{text: "10OwHTs=", key: "0011223344556677"},
			want:    "hello",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESDecrypt(tt.args.text, tt.args.key)
			if !tt.wantErr(t, err, fmt.Sprintf("AESDecrypt(%v, %v)", tt.args.text, tt.args.key)) {
				return
			}
			assert.Equalf(t, tt.want, got, "AESDecrypt(%v, %v)", tt.args.text, tt.args.key)
		})
	}
}
