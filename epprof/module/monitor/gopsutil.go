package monitor

import (
	"os"

	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v4/process"
)

type gopsutil struct {
	p *process.Process
}

func New() Monitor {
	m := gopsutil{}
	err := m.refreshProcess()
	if err != nil {
		elog.Panic("failedToRefreshProcess", elog.FieldErr(err))
	}
	return &m
}

func (m *gopsutil) refreshProcess() error {
	if m.p != nil {
		return nil
	}
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return errors.Wrapf(err, "failed to get process")
	}
	m.p = p
	return nil
}

// ReadMemStats unit is MB
func (m *gopsutil) ReadMemStats() (uint64, float32) {
	mem, err := m.p.MemoryInfo()
	if err != nil {
		if errors.Is(err, process.ErrorProcessNotRunning) {
			err = m.refreshProcess()
			if err != nil {
				elog.Error("failedToRefreshProcess", l.E(err))
				return 0, 0
			}
		}
		elog.Error("failedToReadMemStats", l.E(err))
		return 0, 0
	}
	memPercent, err := m.p.MemoryPercent()
	if err != nil {
		elog.Error("failedToReadMemoryPercentStats", l.E(err))
		return 0, 0
	}
	return mem.RSS / 1024 / 1024, memPercent
}

// ReadCPUStats unit is MB
func (m *gopsutil) ReadCPUStats() float64 {
	stats, err := m.p.CPUPercent()
	if err != nil {
		if errors.Is(err, process.ErrorProcessNotRunning) {
			err = m.refreshProcess()
			if err != nil {
				elog.Error("failedToRefreshProcess", l.E(err))
			}
		}
		return 0
	}
	return stats
}
