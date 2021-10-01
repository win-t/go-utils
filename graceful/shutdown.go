package graceful

import (
	"context"
	"sync"
)

var shutdown struct {
	sync.Mutex
	sync.WaitGroup
}

// ShutdownAndWait cancel the graceful context and
// wait until all registered WaitOnShutdown ctx to completed.
func ShutdownAndWait() {
	shutdown.Lock()
	defer shutdown.Unlock()

	graceful.cancel()
	<-graceful.Done()

	shutdown.Wait()
}

// Register ctx to be waited on shutdown
func WaitOnShutdown(ctx context.Context) {
	shutdown.Lock()
	defer shutdown.Unlock()

	select {
	case <-graceful.Done():
		return
	default:
	}

	select {
	case <-ctx.Done():
		return
	default:
	}

	shutdown.Add(1)
	go func() {
		defer shutdown.Done()
		<-ctx.Done()
	}()
}
