package mineng

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type Loop struct {
	*ECS
	Period time.Duration
}

func New(period time.Duration) *Loop {
	api := NewECS()
	api.Asset(api)
	return &Loop{ECS: api, Period: period}
}

func (loop *Loop) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		t := time.NewTicker(loop.Period)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-t.C:
				loop.ECS.Step()
			}
		}
	})
	return g.Wait()
}
