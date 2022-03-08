package util

import "runtime"

type MemStats struct {
	Sys        float64
	Alloc      float64
	TotalAlloc float64
}

func GetMemoryStatus() MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemStats{
		Sys:        float64(m.Sys) * 1e-6,
		Alloc:      float64(m.Alloc) * 1e-6,
		TotalAlloc: float64(m.TotalAlloc) * 1e-6,
	}
}
