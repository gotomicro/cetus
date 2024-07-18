package epprof

import (
	"encoding/json"
	"time"
)

type options struct {
	cubeOpts *cubeOption
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
		cubeOpts: newMemOptions(),
	}
}

func newMemOptions() *cubeOption {
	return newCubeOpts(defaultMemTriggerValue, defaultMemTriggerPercent, defaultMemTriggerDiff, defaultCPUTriggerPercent, defaultCoolingTime)
}

func WithMemOpts(value, memPercent, memDiff uint64, coolingTime time.Duration) Option {
	return optionFunc(func(opts *options) error {
		opts.cubeOpts.Set(value, memPercent, memDiff, coolingTime)
		return nil
	})
}

func WithCPUOpts(percent uint64) Option {
	return optionFunc(func(opts *options) error {
		opts.cubeOpts.SetCPU(percent)
		return nil
	})
}

type cubeOption struct {
	Enable            bool
	TriggerValue      uint64
	TriggerMemPercent uint64
	TriggerCPUPercent uint64
	TriggerDiff       uint64
	CoolingTime       time.Duration
}

func newCubeOpts(triggerValue, triggerMemPercent, triggerDiff, triggerCPUPercent uint64, coolingTime time.Duration) *cubeOption {
	return &cubeOption{
		Enable:            false,
		TriggerValue:      triggerValue,
		TriggerMemPercent: triggerMemPercent,
		TriggerDiff:       triggerDiff,
		TriggerCPUPercent: triggerCPUPercent,
		CoolingTime:       coolingTime,
	}
}

func (cube *cubeOption) SetCPU(percent uint64) {
	if percent == 0 {
		percent = defaultCPUTriggerPercent
	}
	cube.TriggerCPUPercent = percent
}

func (cube *cubeOption) Set(value, memPercent, memDiff uint64, coolingTime time.Duration) {
	if coolingTime == 0 {
		coolingTime = defaultCoolingTime
	}
	if value == 0 {
		value = defaultMemTriggerValue
	}
	if memDiff == 0 {
		memDiff = defaultMemTriggerDiff
	}
	if memPercent == 0 {
		memPercent = defaultMemTriggerPercent
	}
	cube.TriggerValue = value
	cube.TriggerMemPercent = memPercent
	cube.TriggerDiff = memDiff
	cube.CoolingTime = coolingTime
}

func (cube *cubeOption) String() string {
	tmp, _ := json.Marshal(cube)
	return string(tmp)
}
