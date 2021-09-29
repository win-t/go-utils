// Package graceful.
//
// this package contain context that will be cancelled
// when graceful shutdown is requested
package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

var (
	ctx         context.Context
	cancelFn    context.CancelFunc = func() {}
	wg          *sync.WaitGroup    = &sync.WaitGroup{}
	contextOnce sync.Once
)

// Context for graceful shutdown
func Context() context.Context {
	contextOnce.Do(func() {
		ctx, cancelFn = context.WithCancel(context.Background())
		go func() {
			defer cancelFn()
			c := make(chan os.Signal, 1)
			signal.Notify(c, getInterruptSigs()...)
			select {
			case <-c:
			case <-ctx.Done():
			}
			signal.Stop(c)
		}()
	})
	return ctx
}

// Shutdown cancel the graceful context and wait until the wait counter is zero
func Shutdown() {
	cancelFn()
	wg.Wait()
}

// WaitAdd increase the wait counter
func WaitAdd() {
	wg.Add(1)
}

// WaitDone decrease the wait counter
func WaitDone() {
	wg.Done()
}
