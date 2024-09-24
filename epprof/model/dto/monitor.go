package dto

type CPUInfo struct {
	Frequency uint64  `json:"frequency,omitempty"`
	Quota     float64 `json:"quota,omitempty"`
}
