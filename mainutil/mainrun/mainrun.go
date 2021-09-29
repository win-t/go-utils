// Package mainrun.
package mainrun

import (
	"fmt"
	"os"

	"github.com/payfazz/go-errors/v2"
	"github.com/payfazz/go-errors/v2/trace"

	"github.com/win-t/go-utils/graceful"
)

// Run f, run graceful.Shutdown() after f, and exit with exit code 1 if f returned error.
//
// this function never return.
func Run(f func() error) {
	if err := errors.Catch(func() error {
		defer graceful.Shutdown()
		return f()
	}); err != nil {
		fmt.Fprintln(os.Stderr, errors.FormatWithFilter(err, func(l trace.Location) bool {
			return !l.InPkg("github.com/win-t/go-utils")
		}))
		os.Exit(1)
	}
}
