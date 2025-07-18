package demo

import (
	"context"
	"errors"
	"math/rand"
	"time"
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
	// Simulate variable clone time (500ms to 3s)
	cloneTime := time.Duration(500+rand.Intn(2500)) * time.Millisecond

	// Create a timer for the clone operation
	timer := time.NewTimer(cloneTime)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		// 80% success rate in demo mode
		if rand.Float32() < 0.8 {
			return nil // Success
		}
		// Return a random error
		return DemoErrors[rand.Intn(len(DemoErrors))]
	}
}
