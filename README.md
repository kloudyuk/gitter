# Gitter

A simple utility for testing git server stability by repeatedly cloning repositories.

## Overview

Gitter helps you test the reliability and performance of git servers by continuously cloning repositories while tracking:

- **Success/failure rates** with real-time counters
- **System resources** including goroutines and memory usage (current and peak values)
- **Error history** with timestamps for troubleshooting
- **Runtime duration**

Perfect for monitoring status of git server during changes & upgrades.

## Features

- 🔄 **Continuous Git Cloning** - Repeatedly clone repositories at configurable intervals
- 📊 **Real-time Monitoring** - Live display of success/failure counts and system metrics
- 📈 **Resource Tracking** - Monitor memory usage and goroutines with peak value tracking
- 🚨 **Error Analysis** - Track recent errors with timestamps for troubleshooting
- ⏱️ **Duration Tracking** - See how long the stability test has been running
- 🎨 **Configurable UI** - Adjustable terminal width for different screen sizes
- 🎮 **Demo Mode** - Simulate git operations for testing and demonstration (use `--demo` flag)
- 📝 **Logging** - Errors are logged to files for later analysis (real mode only, not in demo)

## Installation

### Prerequisites

- Go 1.19 or later
- Git (for actual repository cloning)

### Install Directly

```bash
# Install latest version
go install github.com/kloudyuk/gitter@latest

# Install specific tagged version (when available)
go install github.com/kloudyuk/gitter@v0.3.1
```

### Build from Source

```bash
git clone https://github.com/kloudyuk/gitter.git
cd gitter
make build
```

### Check Version

```bash
gitter version
```

## Usage

### Basic Git Server Testing

Test cloning a git repository with default settings (2-second intervals, 10-second timeout):

```bash
gitter clone https://github.com/user/repository.git
```

### Custom Configuration

```bash
# Test with faster intervals and custom timeout
gitter clone https://github.com/user/repo.git --interval 500ms --timeout 30s

# Adjust display width for your terminal
gitter clone https://github.com/user/repo.git --width 150
```

### Demo Mode

Run a simulation without actually cloning repositories (perfect for testing the tool itself):

```bash
# Basic demo mode
gitter clone --demo

# Fast demo with custom settings
gitter clone --demo --interval 500ms --width 100
```

## Command Line Options

### Clone Command

```bash
gitter clone <URL> [flags]
```

**Flags:**

- `-i, --interval duration` - Interval between clones (default: 2s, must be positive)
- `-t, --timeout duration` - Git clone timeout (default: 10s, must be positive)
- `-w, --width int` - Terminal width for display (default: 100, range: 50-300)
- `-e, --error-history int` - Number of recent errors to display (default: 5, must be positive)
- `-d, --demo` - Run in demo mode with simulated git operations

**Note:** When using `--demo` flag, the URL argument becomes optional as the command will use a simulated repository.

### Input Validation

Gitter validates input parameters to ensure reliable operation:

- **Interval**: Must be positive (e.g., `500ms`, `2s`, `1m`)
- **Timeout**: Must be positive (e.g., `10s`, `30s`, `2m`)
- **Width**: Must be between 50 and 300 characters
- **Error History**: Must be positive (e.g., `3`, `10`, `20`)

Invalid inputs will show helpful error messages:

```bash
$ gitter clone --demo --width 30
ERROR: width must be between 50 and 300, got 30

$ gitter clone --demo --error-history 0
ERROR: error-history must be positive, got 0
```

## Display Interface

The live interface shows:

```text
┌────────────────────────────────────────────────────────────────────────────────┐
│                                    Gitter                                      │
│ Config                                                                         │
│ Repo         : https://github.com/user/repo.git                                │
│ Interval     : 2s                                                              │
│ Timeout      : 10s                                                             │
│ Error History: 5                                                               │
│                                                                                │
│ Stats                                                                          │
│ Duration       : 1m30s                                                         │
│ Go Routines    : 5 (max: 8)                                                    │
│ Memory         : 1024 KB (max: 2048 KB)                                        │
│                                                                                │
│ Recent Errors                                                                  │
│ 10s ago: connection timeout                                                    │
│ 45s ago: remote hung up unexpectedly                                           │
│                                                                                │
│ ────────────────────────────────────────────────────────────────────────────── │
│ ⣽ Succeeded: 42                                                                │
│ ⣽ Failed: 3                                                                    │
└────────────────────────────────────────────────────────────────────────────────┘
```

### Key Features Displayed

- **Config Section**: Shows repository URL, interval, timeout, and error history settings
- **Stats Section**: Runtime duration, current/max goroutines, current/max memory usage
- **Recent Errors**: Recent errors with timestamps (configurable history length)
- **Results**: Real-time success/failure counters with animated spinners

## Development

### Complete CI Pipeline

```bash
# Run full CI pipeline: clean, lint, test with coverage, and build
make all
```

### Building

```bash
# Build with automatic version detection
make build

# Build with specific version
VERSION=v1.2.3 make build
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make cover
```

### Linting

```bash
make lint
```

## How It Works

1. **Repository Cloning**: Uses [go-git](https://github.com/go-git/go-git) for efficient in-memory cloning
2. **Concurrency**: Each clone operation runs in a separate goroutine with proper timeout handling
3. **Resource Monitoring**: Tracks system metrics every second using Go's runtime package
4. **Error Tracking**: Maintains a rolling buffer of recent errors with timestamps
5. **UI Updates**: Real-time terminal interface using [Bubble Tea](https://github.com/charmbracelet/bubbletea)

## Output Files

- `gitter.log` - Error log for real clone operations (demo mode creates no log files)
- `coverage.html` - Test coverage report (when running `make cover`)

## Configuration

Gitter uses golangci-lint for code quality. The tool runs with default linting rules to maintain code standards.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linting (`make test lint`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - For the excellent TUI framework
- [go-git](https://github.com/go-git/go-git) - For pure Go git implementation
- [Cobra](https://github.com/spf13/cobra) - For CLI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - For terminal styling

---

**Happy Git Testing!** 🚀
