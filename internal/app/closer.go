package app

import (
	"context"
	"fmt"
	"log/slog"
)

type FuncCloseWithContext func(ctx context.Context) error
type FuncCloseNoContext func() error

type funcClose struct {
	fn  FuncCloseWithContext
	msg string
}

type appCloser struct {
	logger *slog.Logger
	funcs  []funcClose
}

func newCloser(log *slog.Logger) *appCloser {
	if log == nil {
		log = slog.Default()
	}
	return &appCloser{
		funcs:  []funcClose{},
		logger: log.With("component", "logger"),
	}
}

func (c *appCloser) graceFullShutDown(ctx context.Context, fnc funcClose) error {
	chSignal := make(chan struct{}, 1)

	go func() {
		fnc.fn(ctx)
		chSignal <- struct{}{}
		c.logger.Debug(fnc.msg)

	}()

	select {
	case <-chSignal:

		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout shutdown")
	}
}

func (c *appCloser) add(fn FuncCloseNoContext, msg string) {
	if fn == nil {
		return
	}
	fnContext := func(context.Context) error {
		return fn()
	}
	fnc := funcClose{
		msg: msg,
		fn:  fnContext,
	}
	c.funcs = append(c.funcs, fnc)
}

func (c *appCloser) addWithContext(fn FuncCloseWithContext, msg string) {
	if fn == nil {
		return
	}

	fnc := funcClose{
		msg: msg,
		fn:  fn,
	}

	c.funcs = append(c.funcs, fnc)
}

func (c *appCloser) Close(ctx context.Context) error {
	for i := len(c.funcs) - 1; i >= 0; i-- {
		fn := c.funcs[i]
		if err := c.graceFullShutDown(ctx, fn); err != nil {
			return err
		}
	}

	return nil
}
