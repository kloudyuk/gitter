package ui

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestErrorStatsTracking(t *testing.T) {
	// Create error stats with max 3 recent errors
	errorStats := NewErrorStats(3)

	// Add some errors
	err1 := errors.New("first error")
	err2 := errors.New("second error")
	err3 := errors.New("third error")
	err4 := errors.New("fourth error")

	now := time.Now()

	// Add errors
	errorStats.AddError(err1, now)
	errorStats.AddError(err2, now.Add(1*time.Second))
	errorStats.AddError(err3, now.Add(2*time.Second))
	errorStats.AddError(err4, now.Add(3*time.Second))

	// Should have only 3 recent errors (oldest dropped)
	recentErrors := errorStats.GetRecentErrors()
	if len(recentErrors) != 3 {
		t.Errorf("Expected 3 recent errors, got %d", len(recentErrors))
	}

	// Should have total of 4 errors
	if errorStats.GetTotalErrors() != 4 {
		t.Errorf("Expected 4 total errors, got %d", errorStats.GetTotalErrors())
	}

	// Most recent error should be the fourth one
	if recentErrors[len(recentErrors)-1].err.Error() != "fourth error" {
		t.Errorf("Expected most recent error to be 'fourth error', got '%s'",
			recentErrors[len(recentErrors)-1].err.Error())
	}

	// Oldest stored error should be the second one (first was dropped)
	if recentErrors[0].err.Error() != "second error" {
		t.Errorf("Expected oldest stored error to be 'second error', got '%s'",
			recentErrors[0].err.Error())
	}
}

func TestAppStatsMaxTracking(t *testing.T) {
	stats := NewAppStats()

	// Simulate updating goroutines and memory
	stats.UpdateStats(5, 1024)
	if stats.maxGoRoutines != 5 {
		t.Errorf("Expected max goroutines to be 5, got %d", stats.maxGoRoutines)
	}
	if stats.maxMemory != 1024 {
		t.Errorf("Expected max memory to be 1024, got %d", stats.maxMemory)
	}

	// Update with lower values - max should remain
	stats.UpdateStats(3, 512)
	if stats.maxGoRoutines != 5 {
		t.Errorf("Expected max goroutines to remain 5, got %d", stats.maxGoRoutines)
	}
	if stats.maxMemory != 1024 {
		t.Errorf("Expected max memory to remain 1024, got %d", stats.maxMemory)
	}

	// Update with higher values - max should increase
	stats.UpdateStats(8, 2048)
	if stats.maxGoRoutines != 8 {
		t.Errorf("Expected max goroutines to be 8, got %d", stats.maxGoRoutines)
	}
	if stats.maxMemory != 2048 {
		t.Errorf("Expected max memory to be 2048, got %d", stats.maxMemory)
	}
}

// Test only the view generation without UI updates that could hang
func TestConfigView(t *testing.T) {
	styles := NewStyles(100)
	m := model{
		settings: &appSettings{
			repo:         "https://github.com/test/repo.git",
			interval:     2 * time.Second,
			timeout:      10 * time.Second,
			width:        100,
			demoMode:     false,
			errorHistory: 5,
		},
		styles: styles,
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
	if !strings.Contains(configView, "5") {
		t.Error("Config view should contain error history")
	}
}

func TestStatsView(t *testing.T) {
	stats := NewAppStats()
	// Set up test data by directly setting fields for testing
	stats.goRoutines = 5
	stats.maxGoRoutines = 8
	stats.memStats.Alloc = 1024 * 1024                  // 1MB
	stats.maxMemory = 2 * 1024 * 1024                   // 2MB
	stats.startTime = time.Now().Add(-30 * time.Second) // 30 seconds ago

	styles := NewStyles(100)
	m := model{
		stats:  stats,
		styles: styles,
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
