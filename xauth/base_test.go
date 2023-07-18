package xauth

import (
	"testing"
)

func TestParseAppAndSubURL(t *testing.T) {
	type args struct {
		rootURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			args: args{
				rootURL: "",
			},
			want:    "",
			want1:   "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseAppAndSubURL(tt.args.rootURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAppAndSubURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseAppAndSubURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseAppAndSubURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
