package monitor

import (
	"testing"
)

func TestReadMemStats(t *testing.T) {
	tests := []struct {
		name string
		want uint64
	}{
		// TODO: Add test cases.
		{"test1", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ReadMemStats(); got != tt.want {
				t.Errorf("ReadMemStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
