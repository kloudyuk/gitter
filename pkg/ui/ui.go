package ui

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/kloudyuk/gitter/pkg/demo"
	"github.com/kloudyuk/gitter/pkg/git"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Style functions that take width as parameter
func getMainStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(width).
		Padding(0, 1, 0, 1).
		MarginBottom(1)
}

func getTitleStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Bold(true).
		Foreground(lipgloss.Color("#33B5E5")).
		Width(width)
}

func getConfigStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		MarginBottom(1)
}

func getStatsStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func getErrStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func getResultStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width - 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true)
}

type errorInfo struct {
	err       error
	timestamp time.Time
}

type errorStats struct {
	recentErrors []errorInfo
	maxRecent    int
	totalErrors  int
}

type model struct {
	settings   *appSettings
	stats      *appStats
	errorStats *errorStats
	success    result
	fail       result
	resultC    chan error
}

type appSettings struct {
	t        *time.Ticker
	repo     string
	timeout  time.Duration
	interval time.Duration
	log      io.Writer
	demoMode bool
	width    int
}

type appStats struct {
	t             *time.Ticker
	startTime     time.Time
	goRoutines    int
	maxGoRoutines int
	memStats      *runtime.MemStats
	maxMemory     uint64
}

type result struct {
	spinner spinner.Model
	count   int
}

type memStatMsg struct{}
type resultMsg struct {
	err error
}

func updateMemoryStats(c <-chan time.Time) tea.Cmd {
	return func() tea.Msg {
		<-c
		return memStatMsg{}
	}
}

func clone(t <-chan time.Time, repo string, timeout time.Duration, resultC chan<- error) tea.Cmd {
	return func() tea.Msg {
		for {
			<-t
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			select {
			case <-ctx.Done():
				resultC <- ctx.Err()
			default:
				resultC <- git.Clone(ctx, repo)
				cancel()
			}
		}
	}
}

func demoClone(t <-chan time.Time, repo string, timeout time.Duration, resultC chan<- error) tea.Cmd {
	return func() tea.Msg {
		for {
			<-t
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			select {
			case <-ctx.Done():
				resultC <- ctx.Err()
			default:
				resultC <- demo.Clone(ctx, repo)
				cancel()
			}
		}
	}
}

func waitForResults(resultC <-chan error) tea.Cmd {
	return func() tea.Msg {
		return resultMsg{<-resultC}
	}
}

func (m model) Init() tea.Cmd {
	var cloneCmd tea.Cmd
	if m.settings.demoMode {
		cloneCmd = demoClone(m.settings.t.C, m.settings.repo, m.settings.timeout, m.resultC)
	} else {
		cloneCmd = clone(m.settings.t.C, m.settings.repo, m.settings.timeout, m.resultC)
	}

	return tea.Batch(
		m.success.spinner.Tick,
		m.fail.spinner.Tick,
		waitForResults(m.resultC),
		updateMemoryStats(m.stats.t.C),
		cloneCmd,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil
	case memStatMsg:
		runtime.ReadMemStats(m.stats.memStats)
		m.stats.goRoutines = runtime.NumGoroutine()

		// Update max values
		if m.stats.goRoutines > m.stats.maxGoRoutines {
			m.stats.maxGoRoutines = m.stats.goRoutines
		}
		if m.stats.memStats.Alloc > m.stats.maxMemory {
			m.stats.maxMemory = m.stats.memStats.Alloc
		}

		return m, updateMemoryStats(m.stats.t.C)
	case spinner.TickMsg:
		var cmd tea.Cmd
		switch msg.ID {
		case m.success.spinner.ID():
			m.success.spinner, cmd = m.success.spinner.Update(msg)
			return m, cmd
		case m.fail.spinner.ID():
			m.fail.spinner, cmd = m.fail.spinner.Update(msg)
			return m, cmd
		default:
			return m, nil
		}
	case resultMsg:
		if msg.err == nil {
			m.success.count++
		} else {
			m.fail.count++
			m.errorStats.addError(msg.err, time.Now())
			// Only log to file in real mode (not demo mode)
			if m.settings.log != nil {
				_, _ = m.settings.log.Write([]byte(msg.err.Error() + "\n"))
			}
		}
		return m, waitForResults(m.resultC)
	default:
		return m, nil
	}
}

func (m model) View() string {
	width := m.settings.width
	return getMainStyle(width).Render(
		lipgloss.JoinVertical(lipgloss.Top,
			getTitleStyle(width).Render("Gitter"),
			getConfigStyle().Render(m.configView()),
			getStatsStyle().Render(m.statsView()),
			getErrStyle().Render(m.errView()),
			getResultStyle(width).Render(m.resultsView()),
		),
	)
}

func (m model) configView() string {
	return fmt.Sprintf(`%s
Repo     : %s
Interval : %s
Timeout  : %s`,
		title("Config", "#BBBB00"),
		m.settings.repo,
		m.settings.interval,
		m.settings.timeout,
	)
}

func (m model) statsView() string {
	duration := time.Since(m.stats.startTime).Truncate(time.Second)
	return fmt.Sprintf(`%s
Duration       : %s
Go Routines    : %d (max: %d)
Memory         : %d KB (max: %d KB)`,
		title("Stats", "#BBBB00"),
		duration,
		m.stats.goRoutines,
		m.stats.maxGoRoutines,
		(m.stats.memStats.Alloc / 1024),
		(m.stats.maxMemory / 1024),
	)
}

func (m model) resultsView() string {
	return fmt.Sprintf(`%s Succeeded: %d
%s Failed: %d`,
		m.success.spinner.View(), m.success.count,
		m.fail.spinner.View(), m.fail.count,
	)
}

func (m model) errView() string {
	if len(m.errorStats.recentErrors) == 0 {
		return ""
	}

	var errorDisplay []string
	errorDisplay = append(errorDisplay, title("Recent Errors", "#FF0000"))

	// Show recent errors with timestamps
	for i := len(m.errorStats.recentErrors) - 1; i >= 0; i-- {
		errInfo := m.errorStats.recentErrors[i]
		timeAgo := time.Since(errInfo.timestamp).Truncate(time.Second)
		errMsg := errInfo.err.Error()
		if len(errMsg) > 50 {
			errMsg = errMsg[:47] + "..."
		}
		errorDisplay = append(errorDisplay,
			fmt.Sprintf("%s ago: %s",
				timeAgo,
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6666")).Render(errMsg)))
	}

	return "\n" + lipgloss.JoinVertical(lipgloss.Left, errorDisplay...) + "\n"
}

func title(s, color string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Underline(true).
		Render(s)
}

func Start(repo string, interval, timeout time.Duration, width int, demoMode bool, errorHistory int) error {
	var f *os.File
	var err error
	
	// Only create log file for real mode, not demo mode
	if !demoMode {
		f, err = os.Create("gitter.log")
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := f.Close(); closeErr != nil {
				// Log close error, but don't override main error
				_, _ = fmt.Fprintf(os.Stderr, "Warning: failed to close log file: %v\n", closeErr)
			}
		}()
	}
	p := tea.NewProgram(model{
		settings: &appSettings{
			t:        time.NewTicker(interval),
			repo:     repo,
			timeout:  timeout,
			interval: interval,
			log:      f,
			demoMode: demoMode,
			width:    width,
		},
		stats: &appStats{
			t:             time.NewTicker(1 * time.Second),
			startTime:     time.Now(),
			goRoutines:    0,
			maxGoRoutines: 0,
			memStats:      &runtime.MemStats{},
			maxMemory:     0,
		},
		errorStats: &errorStats{
			recentErrors: make([]errorInfo, 0),
			maxRecent:    errorHistory,
			totalErrors:  0,
		},
		success: result{
			spinner: spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")))),
			count:   0,
		},
		fail: result{
			spinner: spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")))),
			count:   0,
		},
		resultC: make(chan error),
	}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

// Helper methods for error stats
func (es *errorStats) addError(err error, timestamp time.Time) {
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

// Helper methods for app stats
func (as *appStats) updateStats(goroutines int, memory uint64) {
	as.goRoutines = goroutines
	if goroutines > as.maxGoRoutines {
		as.maxGoRoutines = goroutines
	}

	as.memStats.Alloc = memory
	if memory > as.maxMemory {
		as.maxMemory = memory
	}
}
