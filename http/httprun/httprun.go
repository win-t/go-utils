// Package httprun.
package httprun

import (
	"context"
	"net/http"

	"github.com/payfazz/go-errors/v2"

	"github.com/win-t/go-utils/http/defserver"
)

func Run(ctx context.Context, s *http.Server) error {
	errCh := make(chan error, 1)
	go func() {
		if s.TLSConfig == nil {
			errCh <- ignoreErrServerClosed(s.ListenAndServe())
		} else {
			errCh <- ignoreErrServerClosed(s.ListenAndServeTLS("", ""))
		}
	}()

	return wait(ctx, errCh, s)
}

func ListenAndServe(ctx context.Context, addr string, handler http.Handler, opts ...defserver.Option) error {
	s, err := defserver.New(addr, handler, opts...)
	if err != nil {
		return errors.Trace(err)
	}
	return Run(ctx, s)
}
