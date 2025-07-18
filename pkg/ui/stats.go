package ui

import (
	"runtime"
	"time"
)

// AppStats tracks application performance statistics
type AppStats struct {
	t             *time.Ticker
	startTime     time.Time
	goRoutines    int
	maxGoRoutines int
	memStats      *runtime.MemStats
	maxMemory     uint64
}

// NewAppStats creates a new AppStats instance
func NewAppStats() *AppStats {
	return &AppStats{
		t:             time.NewTicker(1 * time.Second),
		startTime:     time.Now(),
		goRoutines:    0,
		maxGoRoutines: 0,
		memStats:      &runtime.MemStats{},
		maxMemory:     0,
	}
}

// UpdateStats updates the current stats and tracks maximums
func (as *AppStats) UpdateStats(goroutines int, memory uint64) {
	as.goRoutines = goroutines
	if goroutines > as.maxGoRoutines {
		as.maxGoRoutines = goroutines
	}

	as.memStats.Alloc = memory
	if memory > as.maxMemory {
		as.maxMemory = memory
	}
}

// GetDuration returns the elapsed time since stats tracking started
func (as *AppStats) GetDuration() time.Duration {
	return time.Since(as.startTime).Truncate(time.Second)
}

// GetCurrentMemoryKB returns current memory usage in KB
func (as *AppStats) GetCurrentMemoryKB() uint64 {
	return as.memStats.Alloc / 1024
}

// GetMaxMemoryKB returns maximum memory usage in KB
func (as *AppStats) GetMaxMemoryKB() uint64 {
	return as.maxMemory / 1024
}
