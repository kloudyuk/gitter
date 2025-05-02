package ui

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/kloudyuk/gitter/pkg/git"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var maxWidth = 80

var (
	mainStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Width(maxWidth).
			Padding(0, 1, 0, 1).
			MarginBottom(1)

	titleStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("#33B5E5")).
			Width(maxWidth)

	configStyle = lipgloss.NewStyle().
			MarginBottom(1)

	statsStyle = lipgloss.NewStyle()

	errStyle = lipgloss.NewStyle()

	resultStyle = lipgloss.NewStyle().
			Width(maxWidth - 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true)
)

type model struct {
	settings  *appSettings
	stats     *appStats
	success   result
	fail      result
	lastError error
	resultC   chan error
}

type appSettings struct {
	t        *time.Ticker
	repo     string
	timeout  time.Duration
	interval time.Duration
	log      io.Writer
}

type appStats struct {
	t          *time.Ticker
	goRoutines int
	memStats   *runtime.MemStats
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

func waitForResults(resultC <-chan error) tea.Cmd {
	return func() tea.Msg {
		return resultMsg{<-resultC}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.success.spinner.Tick,
		m.fail.spinner.Tick,
		waitForResults(m.resultC),
		updateMemoryStats(m.stats.t.C),
		clone(m.settings.t.C, m.settings.repo, m.settings.timeout, m.resultC),
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
			m.lastError = msg.err
			m.settings.log.Write([]byte(msg.err.Error() + "\n"))
		}
		return m, waitForResults(m.resultC)
	default:
		return m, nil
	}
}

func (m model) View() string {
	return mainStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			titleStyle.Render("Gitter"),
			configStyle.Render(m.configView()),
			statsStyle.Render(m.statsView()),
			errStyle.Render(m.errView()),
			resultStyle.Render(m.resultsView()),
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
	return fmt.Sprintf(`%s
Go Routines : %d
Memory      : %d KB`,
		title("Stats", "#BBBB00"),
		m.stats.goRoutines,
		(m.stats.memStats.Alloc / 1024),
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
	s := ""
	if m.lastError != nil {
		s = fmt.Sprintf("\nLast Error: %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render(m.lastError.Error()))
	}
	return s
}

func title(s, color string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Underline(true).
		Render(s)
}

func Start(repo string, interval, timeout time.Duration) error {
	f, err := os.Create("gitter.log")
	if err != nil {
		return err
	}
	defer f.Close()
	p := tea.NewProgram(model{
		settings: &appSettings{
			t:        time.NewTicker(interval),
			repo:     repo,
			timeout:  timeout,
			interval: interval,
			log:      f,
		},
		stats: &appStats{
			t:          time.NewTicker(1 * time.Second),
			goRoutines: 0,
			memStats:   &runtime.MemStats{},
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
