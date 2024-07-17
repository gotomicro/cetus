package epprof

import (
	"time"

	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"

	"github.com/gotomicro/cetus/epprof/model/dto"
	"github.com/gotomicro/cetus/epprof/module/forward/webhook"
	"github.com/gotomicro/cetus/epprof/module/monitor"
)

type Epprof struct {
	opts           *options
	memAvg         uint64
	memCoolingTime time.Time

	Forwarder dto.Forwarder
	monitor   monitor.Monitor
}

// New is the entry of epprof
// 1. start monitor
// 2. storage data
// 3. calculate data
// 4. generate pprof files
// 5. upload to oss
func New(opts ...Option) (*Epprof, error) {
	e := &Epprof{
		opts:    newOptions(),
		monitor: monitor.New(),
	}
	for _, o := range opts {
		if err := o.apply(e.opts); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (a *Epprof) Apply(opts ...Option) error {
	for _, o := range opts {
		if err := o.apply(a.opts); err != nil {
			return err
		}
	}
	return nil
}

// EnableMem enables the mem dump.
func (a *Epprof) EnableMem() *Epprof {
	a.opts.memOpts.Enable = true
	return a
}

func (a *Epprof) Start() error {
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

func (a *Epprof) calculateData() {
	logInfo := dto.LogInfo{
		OptAbs:            a.opts.memOpts.TriggerValue,
		OptDiff:           a.opts.memOpts.TriggerDiff,
		OptCoolingTimeSec: int(a.opts.memOpts.CoolingTime.Seconds()),
	}
	if a.isReportCPU(&logInfo) || a.isReportMemory(&logInfo) {
		a.pprof(logInfo)
	}
}

func (a *Epprof) isReportCPU(logInfo *dto.LogInfo) bool {
	usedPercent := a.monitor.ReadCPUStats()
	elog.Debug("cal", l.S("step", "isReportCPU"), l.F64("usedPercent", usedPercent))
	if usedPercent > 0.8 {
		logInfo.CPUPercent = usedPercent
		return true
	}
	return false
}

func (a *Epprof) isReportMemory(logInfo *dto.LogInfo) bool {
	current, usedPercent := a.monitor.ReadMemStats()
	if a.memAvg == 0 {
		a.memAvg = current
		return false
	}
	diff := (float64(current) - float64(a.memAvg)) * 100 / float64(a.memAvg)
	elog.Debug("cal", l.S("step", "isReportMemory"), l.U64("avg", a.memAvg), l.U64("size", current), l.F64("diffPercent", diff), l.F64("usedPercent", float64(usedPercent)), l.A("memOpts", a.opts.memOpts))
	if (current >= a.opts.memOpts.TriggerValue && uint64(diff) >= a.opts.memOpts.TriggerDiff) || usedPercent > 0.8 {
		logInfo.MemoryAbs = current
		logInfo.MemoryDiff = int(diff)
		return true
	}
	a.memAvg = (a.memAvg + current) / 2
	return false
}

func (a *Epprof) pprof(attach dto.LogInfo) {
	if a.memCoolingTime.After(time.Now()) {
		elog.Info("coolingTime")
		return
	}
	a.memCoolingTime = time.Now().Add(a.opts.memOpts.CoolingTime)
	webhook.Webhook(a.Forwarder.WebHook, attach)
}
