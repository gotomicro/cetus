package monitor

type Monitor interface {
	ReadMemStats() (uint64, float32)
	ReadCPUStats() float64
}
