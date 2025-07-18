package ui

import (
	"errors"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestErrorStatsTracking(t *testing.T) {
	// Create error stats with max 3 recent errors
	errorStats := &errorStats{
		recentErrors: make([]errorInfo, 0),
		maxRecent:    3,
		totalErrors:  0,
	}

	// Add some errors
	err1 := errors.New("first error")
	err2 := errors.New("second error")
	err3 := errors.New("third error")
	err4 := errors.New("fourth error")

	now := time.Now()

	// Add errors
	errorStats.addError(err1, now)
	errorStats.addError(err2, now.Add(1*time.Second))
	errorStats.addError(err3, now.Add(2*time.Second))
	errorStats.addError(err4, now.Add(3*time.Second))

	// Should have only 3 recent errors (oldest dropped)
	if len(errorStats.recentErrors) != 3 {
		t.Errorf("Expected 3 recent errors, got %d", len(errorStats.recentErrors))
	}

	// Should have total of 4 errors
	if errorStats.totalErrors != 4 {
		t.Errorf("Expected 4 total errors, got %d", errorStats.totalErrors)
	}

	// Most recent error should be the fourth one
	if errorStats.recentErrors[len(errorStats.recentErrors)-1].err.Error() != "fourth error" {
		t.Errorf("Expected most recent error to be 'fourth error', got '%s'",
			errorStats.recentErrors[len(errorStats.recentErrors)-1].err.Error())
	}

	// Oldest stored error should be the second one (first was dropped)
	if errorStats.recentErrors[0].err.Error() != "second error" {
		t.Errorf("Expected oldest stored error to be 'second error', got '%s'",
			errorStats.recentErrors[0].err.Error())
	}
}

func TestAppStatsMaxTracking(t *testing.T) {
	stats := &appStats{
		startTime:     time.Now(),
		goRoutines:    0,
		maxGoRoutines: 0,
		memStats:      &runtime.MemStats{},
		maxMemory:     0,
	}

	// Simulate updating goroutines and memory
	stats.updateStats(5, 1024)
	if stats.maxGoRoutines != 5 {
		t.Errorf("Expected max goroutines to be 5, got %d", stats.maxGoRoutines)
	}
	if stats.maxMemory != 1024 {
		t.Errorf("Expected max memory to be 1024, got %d", stats.maxMemory)
	}

	// Update with lower values - max should remain
	stats.updateStats(3, 512)
	if stats.maxGoRoutines != 5 {
		t.Errorf("Expected max goroutines to remain 5, got %d", stats.maxGoRoutines)
	}
	if stats.maxMemory != 1024 {
		t.Errorf("Expected max memory to remain 1024, got %d", stats.maxMemory)
	}

	// Update with higher values - max should increase
	stats.updateStats(8, 2048)
	if stats.maxGoRoutines != 8 {
		t.Errorf("Expected max goroutines to be 8, got %d", stats.maxGoRoutines)
	}
	if stats.maxMemory != 2048 {
		t.Errorf("Expected max memory to be 2048, got %d", stats.maxMemory)
	}
}

// Test only the view generation without UI updates that could hang
func TestConfigView(t *testing.T) {
	m := model{
		settings: &appSettings{
			repo:     "https://github.com/test/repo.git",
			interval: 2 * time.Second,
			timeout:  10 * time.Second,
			width:    100,
			demoMode: false,
		},
	}

	configView := m.configView()
	if !strings.Contains(configView, "https://github.com/test/repo.git") {
		t.Error("Config view should contain repository URL")
	}
	if !strings.Contains(configView, "2s") {
		t.Error("Config view should contain interval")
	}
	if !strings.Contains(configView, "10s") {
		t.Error("Config view should contain timeout")
	}
}

func TestStatsView(t *testing.T) {
	m := model{
		stats: &appStats{
			startTime:     time.Now().Add(-30 * time.Second), // 30 seconds ago
			goRoutines:    5,
			maxGoRoutines: 8,
			memStats:      &runtime.MemStats{Alloc: 1024 * 1024}, // 1MB
			maxMemory:     2 * 1024 * 1024,                       // 2MB
		},
	}

	statsView := m.statsView()
	if !strings.Contains(statsView, "5") {
		t.Error("Stats view should contain current goroutines count")
	}
	if !strings.Contains(statsView, "8") {
		t.Error("Stats view should contain max goroutines count")
	}
	if !strings.Contains(statsView, "1024") {
		t.Error("Stats view should contain current memory in KB")
	}
	if !strings.Contains(statsView, "2048") {
		t.Error("Stats view should contain max memory in KB")
	}
}
