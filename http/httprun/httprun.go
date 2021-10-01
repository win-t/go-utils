package httprun

import (
	"context"
	"net/http"
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
