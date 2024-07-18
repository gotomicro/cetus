package monitor

type Monitor interface {
	ReadCPUStats() float64
	ReadMemStats() (uint64, float32)
}
