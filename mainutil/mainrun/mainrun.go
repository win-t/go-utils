// Package mainrun.
package mainrun

import (
	"fmt"
	"os"

	"github.com/payfazz/go-errors/v2"
	"github.com/payfazz/go-errors/v2/trace"

	"github.com/win-t/go-utils/graceful"
)

// Run f
//
// after f is completed, run graceful.ShutdownAndWait
// if f returned error, then run os.Exit(1),
// otherwise run os.Exit(0),
//
// this function never return.
func Run(f func() error) {
	exitCode := 1
	defer func() { os.Exit(exitCode) }()

	err := errors.Catch(func() error {
		defer graceful.ShutdownAndWait()
		return f()
	})

	if err == nil {
		exitCode = 0
		return
	}

	reporter.Lock()
	defer reporter.Unlock()

	if reporter.fn != nil {
		exitCode = reporter.fn(err)
	} else {
		fmt.Fprintln(os.Stderr, errors.FormatWithFilter(err, filterTrace))
	}
}

func filterTrace(l trace.Location) bool {
	return !l.InPkg("github.com/win-t/go-utils")
}
