package epprof

import (
	"time"
)

const (
	defaultCoolingTime = time.Minute

	defaultMemTriggerValue   = 128 // 128MB
	defaultMemTriggerPercent = 80  // 80%
	defaultMemTriggerDiff    = 10  // 10%
)
