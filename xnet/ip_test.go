package xnet

import (
	"testing"
)

func TestGetOutBoundIP(t *testing.T) {
	tests := []struct {
		name    string
		wantIp  string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test-1",
			wantIp:  "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIp, err := GetOutBoundIP()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOutBoundIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIp != tt.wantIp {
				t.Errorf("GetOutBoundIP() gotIp = %v, want %v", gotIp, tt.wantIp)
			}
		})
	}
}
