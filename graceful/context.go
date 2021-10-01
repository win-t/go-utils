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
	setup sync.Once
}

// make sure graceful context is initialized
func init() { Context() }

// Context for graceful shutdown.
func Context() context.Context {
	graceful.setup.Do(func() {
		var cancel context.CancelFunc
		graceful.Context, cancel = context.WithCancel(context.Background())
		go func() {
			defer cancel()
			c := make(chan os.Signal, 1)
			signal.Notify(c, getInterruptSigs()...)
			<-c
			signal.Stop(c)
		}()
	})
	return graceful.Context
}
