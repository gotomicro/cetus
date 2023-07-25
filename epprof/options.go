package epprof

import (
	"time"
)

type options struct {
	memOpts *cubeOption
}

type Option interface {
	apply(*options) error
}

type optionFunc func(*options) error

func (f optionFunc) apply(opts *options) error {
	return f(opts)
}

func newOptions() *options {
	return &options{
		memOpts: newMemOptions(),
	}
}

func newMemOptions() *cubeOption {
	return newCubeOpts(defaultMemTriggerValue, defaultMemTriggerPercent, defaultMemTriggerDiff, defaultCoolingTime)
}

func WithMemOpts(value, percent, diff uint64, coolingTime time.Duration) Option {
	return optionFunc(func(opts *options) error {
		opts.memOpts.Set(value, percent, diff, coolingTime)
		return nil
	})
}

type cubeOption struct {
	Enable         bool
	TriggerValue   uint64
	TriggerPercent uint64
	TriggerDiff    uint64
	CoolingTime    time.Duration
}

func newCubeOpts(triggerValue, triggerPercent, triggerDiff uint64, coolingTime time.Duration) *cubeOption {
	return &cubeOption{
		Enable:         false,
		TriggerValue:   triggerValue,
		TriggerPercent: triggerPercent,
		TriggerDiff:    triggerDiff,
		CoolingTime:    coolingTime,
	}
}

func (cube *cubeOption) Set(value, percent, diff uint64, coolingTime time.Duration) {
	if coolingTime == 0 {
		coolingTime = defaultCoolingTime
	}
	if value == 0 {
		value = defaultMemTriggerValue
	}
	if diff == 0 {
		diff = defaultMemTriggerDiff
	}
	cube.TriggerValue = value
	cube.TriggerPercent = percent
	cube.TriggerDiff = diff
	cube.CoolingTime = coolingTime
}
