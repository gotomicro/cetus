package dto

type Forwarder struct {
	WebHook Webhook `json:"webhook"`
}

type Webhook struct {
	Url     string                 `json:"url"`
	Headers map[string]string      `json:"headers"`
	Body    map[string]interface{} `json:"body"`
}

type AttachInfo struct {
	CurrentAbs  uint64 `json:"currentAbs,omitempty"`
	CurrentDiff int    `json:"currentDiff,omitempty"`

	// Rule
	OptAbs            uint64 `json:"optAbs,omitempty"`
	OptDiff           uint64 `json:"optDiff,omitempty"`
	OptCoolingTimeSec int    `json:"optCoolingTimeSec,omitempty"`
}
