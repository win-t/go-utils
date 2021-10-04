// Package mainrun.
package mainrun

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/payfazz/go-errors/v2"
	"github.com/payfazz/go-errors/v2/trace"
)

// Run f
//
// ctx passed to f will be canceled when graceful shutdown is requested,
// if f returned error or panic, then log it and run os.Exit(1), otherwise run os.Exit(0).
//
// This function never return.
func Run(f func(ctx context.Context) error) {
	exitCode := 1
	defer func() { os.Exit(exitCode) }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer cancel()
		c := make(chan os.Signal, 1)
		signal.Notify(c, getInterruptSigs()...)
		<-c
		signal.Stop(c)
	}()

	err := errors.Catch(func() error { return f(ctx) })
	if err == nil {
		exitCode = 0
		return
	}

	onError.Lock()
	defer onError.Unlock()

	if onError.fn != nil {
		exitCode = onError.fn(err)
		return
	}

	fmt.Fprintln(os.Stderr, errors.FormatWithFilter(err,
		func(l trace.Location) bool { return !l.InPkg("github.com/win-t/go-utils") },
	))
}
