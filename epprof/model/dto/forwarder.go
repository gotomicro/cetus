package dto

type Forwarder struct {
	WebHook Webhook `json:"webhook"`
}

type Webhook struct {
	Url     string                 `json:"url"`
	Headers map[string]string      `json:"headers"`
	Body    map[string]interface{} `json:"body"`
}

type LogInfo struct {
	// cpu
	CPUPercent float64 `json:"cpuPercent,omitempty"`

	// mem
	MemoryAbs  uint64 `json:"memoryAbs,omitempty"`
	MemoryDiff int    `json:"memoryDiff,omitempty"`

	// Rule
	Options string `json:"opts,omitempty"`
}
