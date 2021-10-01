package httprun

import (
	"context"
	"net/http"
	"time"

	"github.com/payfazz/go-errors/v2"

	"github.com/win-t/go-utils/graceful"
	"github.com/win-t/go-utils/http/defserver"
)

func Run(addr string, handler http.HandlerFunc, opts ...defserver.Option) error {
	s, err := defserver.New(addr, handler, opts...)
	if err != nil {
		return errors.Trace(err)
	}

	errCh := make(chan error, 1)
	go func() {
		var err error
		if s.TLSConfig == nil {
			err = s.ListenAndServe()
		} else {
			err = s.ListenAndServeTLS("", "")
		}
		if !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	return wait(errCh, s)
}

func wait(errCh chan error, s *http.Server) error {
	select {
	case err := <-errCh:
		return errors.Trace(err)
	case <-graceful.Context().Done():
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		return errors.New("cannot gracefuly shutdown server in 5 second")
	}

	return nil
}
