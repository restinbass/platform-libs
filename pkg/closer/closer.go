package closer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"
)

const shutdownTimeout = 5 * time.Second

type (
	// Logger -
	Logger interface {
		Info(ctx context.Context, msg string, fields ...zap.Field)
		Error(ctx context.Context, msg string, fields ...zap.Field)
	}

	// Closer -
	Closer struct {
		mu     sync.Mutex
		once   sync.Once
		done   chan struct{}
		funcs  []func(context.Context) error
		logger Logger
	}
)

var (
	globalCloser = New(&noopLogger{})
)

// AddNamed -
func AddNamed(name string, f func(context.Context) error) {
	globalCloser.AddNamed(name, f)
}

// CloseAll -
func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

// Configure -
func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

// New -
func New(logger Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}

	return c
}

// handleSignals -
func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		ctx := context.Background()
		c.logger.Info(ctx, "got os.Signal, starting graceful shutdown...")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, shutdownTimeout)
		defer shutdownCancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error(ctx, "failed to clouse resources: %v", zap.Error(err))
		}
	case <-c.done:
	}
}

// AddNamed -
func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.logger.Info(ctx, fmt.Sprintf("closing: %s...", name))

		err := f(ctx)
		duration := time.Since(start)
		if err != nil {
			c.logger.Error(ctx, fmt.Sprintf("failed to close %s, err: %v (taken: %s)", name, err, duration))
		} else {
			c.logger.Info(ctx, fmt.Sprintf("resource %s successfully closed (taken: %s)", name, duration))
		}

		return err
	})
}

// Add -
func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f...)
}

// CloseAll -
func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.logger.Info(ctx, "no functions to close")
			return
		}

		c.logger.Info(ctx, "starting graceful shutdown process...")

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]
			wg.Add(1)
			go func(f func(context.Context) error) {
				defer wg.Done()

				// Защита от паники
				defer func() {
					if r := recover(); r != nil {
						errCh <- errors.New("panic recovered in closer")
						c.logger.Error(ctx, "panic consumed in Close() function", zap.Any("error", r))
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(f)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for {
			select {
			case <-ctx.Done():
				c.logger.Info(ctx, "context cancelled during closing...", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				if !ok {
					c.logger.Info(ctx, "all resources closed successfully")
					return
				}
				c.logger.Error(ctx, "error when closing some resources", zap.Error(err))
				if result == nil {
					result = err
				}
			}
		}
	})

	return result
}
