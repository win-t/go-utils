// Package gracefulhttp.
//
// This package provide utility to graceful shutdown http.Server.
package gracefulhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/payfazz/go-errors/v2"

	"github.com/win-t/go-utils/graceful"
)

// Shutdown will create goroutine to shutdown s.
func Shutdown(s *http.Server, timeout time.Duration) func() error {
	var doneErr error
	doneCh := make(chan struct{})
	graceful.WaitAdd()
	errors.Go(func(err error) {
		doneErr = err
		close(doneCh)
	}, func() error {
		defer graceful.WaitDone()
		<-graceful.Context().Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := s.Shutdown(shutdownCtx); err != nil {
			return errors.Trace(err)
		}
		return nil
	})
	return func() error {
		<-doneCh
		return doneErr
	}
}
