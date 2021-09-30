// Package graceful.
//
// This package contain context that will be cancelled
// when graceful shutdown is requested.
//
// On POSIX system, this means that when SIGINT or SIGTERM
// is caught.
package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

var graceful struct {
	context.Context
	cancel context.CancelFunc
	once   sync.Once
}

// Context for graceful shutdown.
func Context() context.Context {
	graceful.once.Do(func() {
		graceful.Context, graceful.cancel = context.WithCancel(context.Background())
		go func() {
			defer graceful.cancel()
			c := make(chan os.Signal, 1)
			signal.Notify(c, getInterruptSigs()...)
			select {
			case <-c:
			case <-graceful.Done():
			}
			signal.Stop(c)
		}()
	})
	return graceful.Context
}
