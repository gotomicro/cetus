package xcheck

import (
	"time"

	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"

	"github.com/gotomicro/cetus/xcheck/model/dto"
	"github.com/gotomicro/cetus/xcheck/module/forward/webhook"
	"github.com/gotomicro/cetus/xcheck/module/monitor"
)

type xcheck struct {
	opts           *options
	memAvg         uint64
	memCoolingTime time.Time

	Forwarder dto.Forwarder
}

// New is the entry of xcheck
// 1. start monitor
// 2. storage data
// 3. calculate data
// 4. generate pprof files
// 5. upload to oss
func New(opts ...Option) (*xcheck, error) {
	xcheck := &xcheck{
		opts: newOptions(),
	}
	for _, o := range opts {
		if err := o.apply(xcheck.opts); err != nil {
			return nil, err
		}
	}
	return xcheck, nil
}

func (a *xcheck) Apply(opts ...Option) error {
	for _, o := range opts {
		if err := o.apply(a.opts); err != nil {
			return err
		}
	}
	return nil
}

// EnableMem enables the mem dump.
func (a *xcheck) EnableMem() *xcheck {
	a.opts.memOpts.Enable = true
	return a
}

func (a *xcheck) Start() error {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				a.calculateData()
			}
		}
	}()
	return nil
}

func (a *xcheck) calculateData() {
	current, usedPercent := monitor.ReadMemStats()

	if a.memAvg == 0 {
		a.memAvg = current
		return
	}
	diff := (float64(current) - float64(a.memAvg)) * 100 / float64(a.memAvg)
	elog.Debug("cal", l.UI64("avg", a.memAvg), l.UI64("size", current), l.F64("diffPercent", diff), l.F64("usedPercent", usedPercent), l.A("memOpts", a.opts.memOpts))
	if current >= a.opts.memOpts.TriggerValue && uint64(diff) >= a.opts.memOpts.TriggerDiff {
		a.pprof(dto.AttachInfo{
			CurrentAbs:        current,
			CurrentDiff:       int(diff),
			OptAbs:            a.opts.memOpts.TriggerValue,
			OptDiff:           a.opts.memOpts.TriggerDiff,
			OptCoolingTimeSec: int(a.opts.memOpts.CoolingTime.Seconds()),
		})
	}
	a.memAvg = (a.memAvg + current) / 2
}

func (a *xcheck) pprof(attach dto.AttachInfo) {
	if a.memCoolingTime.After(time.Now()) {
		elog.Info("coolingTime")
		return
	}
	a.memCoolingTime = time.Now().Add(a.opts.memOpts.CoolingTime)
	webhook.Webhook(a.Forwarder.WebHook, attach)
}
