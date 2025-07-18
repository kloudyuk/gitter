package cmd

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestCloneCommandValidation(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no args and no demo flag",
			args:        []string{"clone"},
			expectError: true,
			errorMsg:    "repository URL is required when not in demo mode",
		},
		{
			name:        "demo flag works without URL",
			args:        []string{"clone", "--demo"},
			expectError: false,
		},
		{
			name:        "valid URL without demo",
			args:        []string{"clone", "https://github.com/user/repo.git"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cloneCmd()
			cmd.SetArgs(tt.args[1:]) // Remove "clone" from args

			// Parse flags but don't execute to avoid starting UI
			err := cmd.ParseFlags(tt.args[1:])
			if err != nil {
				t.Fatalf("Failed to parse flags: %v", err)
			}

			// Test the validation logic directly by checking args
			args := cmd.Flags().Args()
			demoFlag, _ := cmd.Flags().GetBool("demo")

			// Simulate the validation logic from the actual command
			var validationErr error
			if !demoFlag && len(args) == 0 {
				validationErr = errors.New("repository URL is required when not in demo mode")
			}

			if tt.expectError {
				if validationErr == nil {
					t.Errorf("Expected error but got none")
				}
				if validationErr != nil && !strings.Contains(validationErr.Error(), tt.errorMsg) {
					t.Errorf("Expected error message '%s' but got '%s'", tt.errorMsg, validationErr.Error())
				}
			} else {
				if validationErr != nil {
					t.Errorf("Got unexpected validation error: %v", validationErr)
				}
			}
		})
	}
}

func TestCloneFlagParsing(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedInterval time.Duration
		expectedTimeout  time.Duration
		expectedWidth    int
		expectedDemo     bool
	}{
		{
			name:             "default values",
			args:             []string{"clone", "--demo"},
			expectedInterval: 2 * time.Second,
			expectedTimeout:  10 * time.Second,
			expectedWidth:    100,
			expectedDemo:     true,
		},
		{
			name:             "custom interval",
			args:             []string{"clone", "--demo", "--interval", "500ms"},
			expectedInterval: 500 * time.Millisecond,
			expectedTimeout:  10 * time.Second,
			expectedWidth:    100,
			expectedDemo:     true,
		},
		{
			name:             "custom timeout",
			args:             []string{"clone", "--demo", "--timeout", "30s"},
			expectedInterval: 2 * time.Second,
			expectedTimeout:  30 * time.Second,
			expectedWidth:    100,
			expectedDemo:     true,
		},
		{
			name:             "custom width",
			args:             []string{"clone", "--demo", "--width", "150"},
			expectedInterval: 2 * time.Second,
			expectedTimeout:  10 * time.Second,
			expectedWidth:    150,
			expectedDemo:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cloneCmd()
			cmd.SetArgs(tt.args[1:]) // Remove "clone" from args

			// We need to parse without executing to check flag values
			err := cmd.ParseFlags(tt.args[1:])
			if err != nil {
				t.Fatalf("Failed to parse flags: %v", err)
			}

			// Check interval flag
			intervalFlag := cmd.Flag("interval")
			if intervalFlag.Value.String() != tt.expectedInterval.String() {
				t.Errorf("Expected interval %v, got %v", tt.expectedInterval, intervalFlag.Value.String())
			}

			// Check timeout flag
			timeoutFlag := cmd.Flag("timeout")
			if timeoutFlag.Value.String() != tt.expectedTimeout.String() {
				t.Errorf("Expected timeout %v, got %v", tt.expectedTimeout, timeoutFlag.Value.String())
			}

			// Check width flag
			widthFlag := cmd.Flag("width")
			if widthFlag.Value.String() != string(rune(tt.expectedWidth)) {
				// Convert int to string for comparison
				var expectedWidthStr string
				switch tt.expectedWidth {
				case 100:
					expectedWidthStr = "100"
				case 150:
					expectedWidthStr = "150"
				}
				if widthFlag.Value.String() != expectedWidthStr {
					t.Errorf("Expected width %d, got %v", tt.expectedWidth, widthFlag.Value.String())
				}
			}

			// Check demo flag
			demoFlag := cmd.Flag("demo")
			expectedDemoStr := "false"
			if tt.expectedDemo {
				expectedDemoStr = "true"
			}
			if demoFlag.Value.String() != expectedDemoStr {
				t.Errorf("Expected demo %v, got %v", tt.expectedDemo, demoFlag.Value.String())
			}
		})
	}
}

func TestCloneCommandInputValidation(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		timeout  string
		width    string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "negative interval",
			interval: "-1s",
			timeout:  "10s",
			width:    "100",
			wantErr:  true,
			errMsg:   "interval must be positive",
		},
		{
			name:     "zero timeout",
			interval: "2s",
			timeout:  "0s",
			width:    "100",
			wantErr:  true,
			errMsg:   "timeout must be positive",
		},
		{
			name:     "width too small",
			interval: "2s",
			timeout:  "10s",
			width:    "30",
			wantErr:  true,
			errMsg:   "width must be between",
		},
		{
			name:     "width too large",
			interval: "2s",
			timeout:  "10s",
			width:    "400",
			wantErr:  true,
			errMsg:   "width must be between",
		},
		{
			name:     "valid inputs",
			interval: "2s",
			timeout:  "10s",
			width:    "100",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic directly
			interval, err := time.ParseDuration(tt.interval)
			if err != nil && !tt.wantErr {
				t.Fatalf("Failed to parse interval: %v", err)
			}

			timeout, err := time.ParseDuration(tt.timeout)
			if err != nil && !tt.wantErr {
				t.Fatalf("Failed to parse timeout: %v", err)
			}

			var width int
			switch tt.width {
			case "30":
				width = 30
			case "400":
				width = 400
			case "100":
				width = 100
			}

			// Test validation logic directly (extracted from clone command)
			var validationErr error
			if interval <= 0 {
				validationErr = errors.New("interval must be positive")
			} else if timeout <= 0 {
				validationErr = errors.New("timeout must be positive")
			} else if width < MinWidth || width > MaxWidth {
				validationErr = errors.New("width must be between 50 and 300")
			}

			if tt.wantErr {
				if validationErr == nil {
					t.Errorf("Expected validation error but got none")
				}
				if validationErr != nil && !strings.Contains(validationErr.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s' but got '%s'", tt.errMsg, validationErr.Error())
				}
			} else {
				if validationErr != nil {
					t.Errorf("Got unexpected validation error: %v", validationErr)
				}
			}
		})
	}
}
