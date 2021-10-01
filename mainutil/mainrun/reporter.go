package mainrun

import "sync"

var reporter struct {
	sync.Mutex
	fn func(error) int
}

func SetReporter(fn func(error) int) {
	reporter.Lock()
	defer reporter.Unlock()

	reporter.fn = fn
}
