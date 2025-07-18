package demo

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"
)

const (
	// Demo timing constants
	MinCloneTime   = 500 * time.Millisecond // Minimum simulated clone time
	MaxCloneTime   = 3 * time.Second        // Maximum simulated clone time
	SuccessRate    = 0.8                    // 80% success rate in demo mode
	CloneTimeRange = 2500                   // Range in milliseconds (MaxCloneTime - MinCloneTime)
)

// DemoErrors represents various types of errors that can occur during git operations
var DemoErrors = []error{
	errors.New("connection timeout"),
	errors.New("repository not found"),
	errors.New("authentication failed"),
	errors.New("network unreachable"),
	errors.New("remote hung up unexpectedly"),
	errors.New("permission denied"),
	errors.New("could not resolve host"),
	errors.New("ssl certificate problem"),
}

// Clone simulates a git clone operation with realistic timing and error patterns
func Clone(ctx context.Context, repo string) error {
	// Simulate variable clone time using constants
	cloneTime := MinCloneTime + time.Duration(rand.IntN(CloneTimeRange))*time.Millisecond

	// Create a timer for the clone operation
	timer := time.NewTimer(cloneTime)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		// Use SuccessRate constant
		if rand.Float32() < SuccessRate {
			return nil // Success
		}
		// Return a random error
		return DemoErrors[rand.IntN(len(DemoErrors))]
	}
}
