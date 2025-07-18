package cmd

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	// Test with default version
	originalVersion := Version
	defer func() { Version = originalVersion }()

	Version = "dev"
	version := getVersion()
	if version == "" {
		t.Error("getVersion() should never return empty string")
	}

	// Test with custom version
	Version = "v1.2.3"
	version = getVersion()
	if version != "v1.2.3" {
		t.Errorf("Expected v1.2.3, got %s", version)
	}
}
