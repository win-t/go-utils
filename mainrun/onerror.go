package mainrun

import "sync"

var onError struct {
	sync.Mutex
	fn func(error) int
}

// When function that passed to Run is returned error or panic,
// run f, the returned int will be used to os.Exit function.
func OnError(f func(error) int) {
	onError.Lock()
	defer onError.Unlock()

	onError.fn = f
}
