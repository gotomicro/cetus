package monitor

import (
	"os"

	"github.com/shirou/gopsutil/v3/process"

	"github.com/gotomicro/cetus/xcheck/model/dto"
)

func New(m dto.Monitor) {

}

// ReadMemStats unit is MB
func ReadMemStats() (uint64, float64) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, 0
	}
	mem, err := p.MemoryInfo()
	if err != nil {
		return 0, 0
	}
	// fmt.Printf("TotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	// usedPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(mem.RSS)/float64(m.Sys)), 64)
	return mem.RSS / 1024 / 1024, 0
}
