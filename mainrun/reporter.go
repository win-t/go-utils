package mainrun

import "sync"

var reporter struct {
	sync.Mutex
	fn func(error) int
}

// Call reporterFunc if function passed to Run returning error or panic.
// reporterFunc return exit code of the process.
func SetReporter(reporterFunc func(error) int) {
	reporter.Lock()
	defer reporter.Unlock()

	reporter.fn = reporterFunc
}
