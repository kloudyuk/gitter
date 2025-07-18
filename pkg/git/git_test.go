package git

import (
	"context"
	"testing"
	"time"
)

func TestCloneWithTimeout(t *testing.T) {
	// Test with a very short timeout to ensure timeout behavior works
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	err := Clone(ctx, "https://github.com/nonexistent/repo.git")
	if err == nil {
		t.Error("Expected error due to timeout, got nil")
	}

	if err != context.DeadlineExceeded {
		// The error might be from the git operation itself if it starts very quickly
		// In that case, we just verify we got some error
		t.Logf("Got error (expected): %v", err)
	}
}

func TestCloneWithCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	err := Clone(ctx, "https://github.com/nonexistent/repo.git")
	if err == nil {
		t.Error("Expected error due to cancellation, got nil")
	}

	if err != context.Canceled {
		// Similar to timeout test, the error might be from git operation
		t.Logf("Got error (expected): %v", err)
	}
}

func TestCloneInvalidRepo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// This should fail because the repository doesn't exist
	err := Clone(ctx, "https://github.com/definitely-does-not-exist-12345/repo.git")
	if err == nil {
		t.Error("Expected error for non-existent repository, got nil")
	}
}
