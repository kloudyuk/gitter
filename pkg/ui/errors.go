package ui

import (
	"time"
)

// ErrorStats tracks error statistics with configurable history
type ErrorStats struct {
	recentErrors []errorInfo
	maxRecent    int
	totalErrors  int
}

// NewErrorStats creates a new ErrorStats instance
func NewErrorStats(maxRecent int) *ErrorStats {
	return &ErrorStats{
		recentErrors: make([]errorInfo, 0),
		maxRecent:    maxRecent,
		totalErrors:  0,
	}
}

// AddError adds an error to the tracking system
func (es *ErrorStats) AddError(err error, timestamp time.Time) {
	es.totalErrors++
	errorInfo := errorInfo{
		err:       err,
		timestamp: timestamp,
	}
	es.recentErrors = append(es.recentErrors, errorInfo)

	// Keep only the most recent errors
	if len(es.recentErrors) > es.maxRecent {
		es.recentErrors = es.recentErrors[1:]
	}
}

// GetRecentErrors returns the recent errors for display
func (es *ErrorStats) GetRecentErrors() []errorInfo {
	return es.recentErrors
}

// GetTotalErrors returns the total number of errors encountered
func (es *ErrorStats) GetTotalErrors() int {
	return es.totalErrors
}
