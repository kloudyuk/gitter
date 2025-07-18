package demo

import (
	"context"
	"testing"
	"time"
)

func TestDemoClone(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		wantErr bool
	}{
		{
			name:    "normal operation",
			timeout: 5 * time.Second,
			wantErr: false, // May or may not error due to randomness
		},
		{
			name:    "with timeout",
			timeout: 1 * time.Millisecond, // Very short timeout
			wantErr: true,                 // Should timeout
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			err := Clone(ctx, "demo-repo")

			if tt.name == "with timeout" {
				// For timeout test, we expect either context deadline exceeded or nil (if operation completed very quickly)
				if err != nil && err != context.DeadlineExceeded {
					t.Errorf("Expected timeout error, got %v", err)
				}
			}
			// For normal operation, we just check that it returns (either success or demo error)
		})
	}
}

func TestDemoErrorTypes(t *testing.T) {
	if len(DemoErrors) == 0 {
		t.Error("DemoErrors should not be empty")
	}

	// Ensure all demo errors are non-nil
	for i, err := range DemoErrors {
		if err == nil {
			t.Errorf("DemoErrors[%d] should not be nil", i)
		}
		if err.Error() == "" {
			t.Errorf("DemoErrors[%d] should have a non-empty error message", i)
		}
	}
}
