package ui

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// errorInfo represents an error with its timestamp
type errorInfo struct {
	err       error
	timestamp time.Time
}

type memStatMsg struct{}

type resultMsg struct {
	err error
}

type model struct {
	settings    *appSettings
	stats       *AppStats
	errorStats  *ErrorStats
	success     result
	fail        result
	resultC     chan error
	cloneRunner *CloneRunner
	styles      *Styles
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

type result struct {
	spinner spinner.Model
	count   int
}

func updateMemoryStats(c <-chan time.Time) tea.Cmd {
	return func() tea.Msg {
		<-c
		return memStatMsg{}
	}
}

func waitForResults(resultC <-chan error) tea.Cmd {
	return func() tea.Msg {
		return resultMsg{<-resultC}
	}
}

func (m model) Init() tea.Cmd {
	cloneCmd := m.cloneRunner.Start()

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
		goroutines := runtime.NumGoroutine()
		memory := m.stats.memStats.Alloc

		// Update stats using the new method
		m.stats.UpdateStats(goroutines, memory)

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
			m.errorStats.AddError(msg.err, time.Now())
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
	return m.styles.Main().Render(
		lipgloss.JoinVertical(lipgloss.Top,
			m.styles.Title().Render("Gitter"),
			m.styles.Config().Render(m.configView()),
			m.styles.Stats().Render(m.statsView()),
			m.styles.Error().Render(m.errView()),
			m.styles.Result().Render(m.resultsView()),
		),
	)
}

func (m model) configView() string {
	return fmt.Sprintf(`%s
Repo     : %s
Interval : %s
Timeout  : %s`,
		m.styles.SectionTitle("Config", "#BBBB00"),
		m.settings.repo,
		m.settings.interval,
		m.settings.timeout,
	)
}

func (m model) statsView() string {
	duration := m.stats.GetDuration()
	return fmt.Sprintf(`%s
Duration       : %s
Go Routines    : %d (max: %d)
Memory         : %d KB (max: %d KB)`,
		m.styles.SectionTitle("Stats", "#BBBB00"),
		duration,
		m.stats.goRoutines,
		m.stats.maxGoRoutines,
		m.stats.GetCurrentMemoryKB(),
		m.stats.GetMaxMemoryKB(),
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
	recentErrors := m.errorStats.GetRecentErrors()
	if len(recentErrors) == 0 {
		return ""
	}

	var errorDisplay []string
	errorDisplay = append(errorDisplay, m.styles.SectionTitle("Recent Errors", "#FF0000"))

	// Show recent errors with timestamps
	for i := len(recentErrors) - 1; i >= 0; i-- {
		errInfo := recentErrors[i]
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

	// Create the channel for results
	resultC := make(chan error)

	// Create new components using constructors
	stats := NewAppStats()
	errorStats := NewErrorStats(errorHistory)
	styles := NewStyles(width)
	cloneRunner := NewCloneRunner(demoMode, time.NewTicker(interval).C, repo, timeout, resultC)

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
		stats:       stats,
		errorStats:  errorStats,
		styles:      styles,
		cloneRunner: cloneRunner,
		success: result{
			spinner: spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")))),
			count:   0,
		},
		fail: result{
			spinner: spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")))),
			count:   0,
		},
		resultC: resultC,
	}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
