package graceful

import "sync"

var shutdown struct {
	sync.Mutex
	fn []func()
}

// Shutdown cancel the graceful context and
// wait until all registered function to completed.
func Shutdown() {
	shutdown.Lock()
	defer shutdown.Unlock()

	Context() // make sure graceful context is initialized
	graceful.cancel()

	var wg sync.WaitGroup
	for _, fn := range shutdown.fn {
		fn := fn
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn()
		}()
	}
	wg.Wait()

	shutdown.fn = nil
}

// Register function f to be called when Shutdown is called,
// if graceful context is already done then f is called immediately.
func RegisterOnShutdown(f func()) {
	shutdown.Lock()
	defer shutdown.Unlock()

	select {
	case <-graceful.Done():
		go f()
	default:
		shutdown.fn = append(shutdown.fn, f)
	}
}
