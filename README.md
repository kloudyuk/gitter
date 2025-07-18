# Gitter

A simple utility for testing git server stability by repeatedly cloning repositories and monitoring system resources.

## Overview

Gitter helps you test the reliability and performance of git servers by continuously cloning repositories while tracking:

- **Success/failure rates** with real-time counters
- **System resources** including goroutines and memory usage (current and peak values)
- **Error history** with timestamps and error rate calculations
- **Runtime duration** to monitor long-term stability

Perfect for testing git server upgrades, network stability, or repository access patterns.

## Features

- 🔄 **Continuous Git Cloning** - Repeatedly clone repositories at configurable intervals
- 📊 **Real-time Monitoring** - Live display of success/failure counts and system metrics
- 📈 **Resource Tracking** - Monitor memory usage and goroutines with peak value tracking
- 🚨 **Error Analysis** - Track recent errors with timestamps and calculate error rates
- ⏱️ **Duration Tracking** - See how long the stability test has been running
- 🎨 **Configurable UI** - Adjustable terminal width for different screen sizes
- 🎮 **Demo Mode** - Simulate git operations for testing and demonstration
- 📝 **Logging** - All errors are logged to files for later analysis

## Installation

### Prerequisites

- Go 1.19 or later
- Git (for actual repository cloning)

### Build from Source

```bash
git clone https://github.com/kloudyuk/gitter.git
cd gitter
make build
```

### Install Directly

```bash
go install github.com/kloudyuk/gitter@latest
```

## Usage

### Basic Git Server Testing

Test a git repository with default settings (2-second intervals, 10-second timeout):

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
# Basic demo
gitter demo

# Fast demo with custom settings
gitter demo --interval 500ms --width 100
```

## Command Line Options

### Clone Command

```bash
gitter clone <URL> [flags]
```

**Flags:**

- `-i, --interval duration` - Interval between clones (default: 2s)
- `-t, --timeout duration` - Git clone timeout (default: 10s)
- `-w, --width int` - Terminal width for display (default: 120)

### Demo Command

```bash
gitter demo [flags]
```

**Flags:**

- `-i, --interval duration` - Interval between simulated clones (default: 2s)
- `-t, --timeout duration` - Simulated git clone timeout (default: 10s)
- `-w, --width int` - Terminal width for display (default: 120)

## Display Interface

The live interface shows:

```text
┌────────────────────────────────────────────────────────────────────────────────┐
│                                    Gitter                                      │
│ Config                                                                         │
│ Repo     : https://github.com/user/repo.git                                   │
│ Interval : 2s                                                                  │
│ Timeout  : 10s                                                                 │
│                                                                                │
│ Stats                                                                          │
│ Duration       : 1m30s                                                         │
│ Go Routines    : 5 (max: 8)                                                    │
│ Memory         : 1024 KB (max: 2048 KB)                                        │
│                                                                                │
│ Recent Errors                                                                  │
│ Error Rate: 2.5/min                                                            │
│ 10s ago: connection timeout                                                    │
│ 45s ago: remote hung up unexpectedly                                           │
│                                                                                │
│ ────────────────────────────────────────────────────────────────────────────── │
│ ✓ Succeeded: 42                                                                │
│ ✗ Failed: 3                                                                    │
└────────────────────────────────────────────────────────────────────────────────┘
```

### Key Features Displayed

- **Config Section**: Shows repository URL, interval, and timeout settings
- **Stats Section**: Runtime duration, current/max goroutines, current/max memory usage
- **Recent Errors**: Last 5 errors with timestamps and error rate per minute
- **Results**: Real-time success/failure counters with animated spinners

## Development

### Building

```bash
make build
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

### Cleaning

```bash
make clean
```

## How It Works

1. **Repository Cloning**: Uses [go-git](https://github.com/go-git/go-git) for efficient in-memory cloning
2. **Concurrency**: Each clone operation runs in a separate goroutine with proper timeout handling
3. **Resource Monitoring**: Tracks system metrics every second using Go's runtime package
4. **Error Tracking**: Maintains a rolling buffer of recent errors with timestamps
5. **UI Updates**: Real-time terminal interface using [Bubble Tea](https://github.com/charmbracelet/bubbletea)

## Use Cases

### Git Server Testing

- Test server stability during upgrades
- Validate load balancing configurations
- Monitor repository access patterns
- Stress test authentication systems

### Network Reliability

- Test git operations over unstable connections
- Monitor clone performance across different networks
- Validate VPN or proxy configurations

### Performance Monitoring

- Baseline git server performance
- Monitor resource usage during peak loads
- Track long-term stability trends

## Output Files

- `gitter.log` - Error log for clone operations
- `gitter-demo.log` - Error log for demo mode
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
