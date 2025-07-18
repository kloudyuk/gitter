package git

import (
	"context"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

func Clone(ctx context.Context, repo string) error {
	errC := make(chan error, 1)
	go func() {
		_, err := gogit.CloneContext(ctx, memory.NewStorage(), nil, &gogit.CloneOptions{
			URL:          repo,
			SingleBranch: true,
			NoCheckout:   true,
			Depth:        1,
		})
		errC <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errC:
		return err
	}
}
