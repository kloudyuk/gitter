package ui

import (
	"context"
	"time"

	"github.com/kloudyuk/gitter/pkg/demo"
	"github.com/kloudyuk/gitter/pkg/git"

	tea "github.com/charmbracelet/bubbletea"
)

// CloneOperation defines the interface for clone operations
type CloneOperation interface {
	Execute(ctx context.Context, repo string) error
}

// RealCloneOperation implements actual git cloning
type RealCloneOperation struct{}

func (r *RealCloneOperation) Execute(ctx context.Context, repo string) error {
	return git.Clone(ctx, repo)
}

// DemoCloneOperation implements simulated git cloning
type DemoCloneOperation struct{}

func (d *DemoCloneOperation) Execute(ctx context.Context, repo string) error {
	return demo.Clone(ctx, repo)
}

// CloneRunner handles the execution of clone operations with timing
type CloneRunner struct {
	operation CloneOperation
	ticker    <-chan time.Time
	repo      string
	timeout   time.Duration
	resultC   chan<- error
}

// NewCloneRunner creates a new CloneRunner
func NewCloneRunner(demoMode bool, ticker <-chan time.Time, repo string, timeout time.Duration, resultC chan<- error) *CloneRunner {
	var operation CloneOperation
	if demoMode {
		operation = &DemoCloneOperation{}
	} else {
		operation = &RealCloneOperation{}
	}

	return &CloneRunner{
		operation: operation,
		ticker:    ticker,
		repo:      repo,
		timeout:   timeout,
		resultC:   resultC,
	}
}

// Start returns a tea.Cmd that runs the clone operations
func (cr *CloneRunner) Start() tea.Cmd {
	return func() tea.Msg {
		for {
			<-cr.ticker
			ctx, cancel := context.WithTimeout(context.Background(), cr.timeout)
			select {
			case <-ctx.Done():
				cr.resultC <- ctx.Err()
			default:
				cr.resultC <- cr.operation.Execute(ctx, cr.repo)
				cancel()
			}
		}
	}
}
