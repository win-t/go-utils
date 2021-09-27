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
	contextOnce sync.Once
)

// Context for graceful shutdown
func Context() context.Context {
	contextOnce.Do(func() {
		var cancelFn context.CancelFunc
		ctx, cancelFn = context.WithCancel(context.Background())
		go func() {
			defer cancelFn()
			c := make(chan os.Signal, 1)
			signal.Notify(c, getInterruptSigs()...)
			<-c
			signal.Stop(c)
		}()
	})
	return ctx
}
