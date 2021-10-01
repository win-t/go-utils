// Package mainhttp.
package mainhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/payfazz/go-errors/v2"

	"github.com/win-t/go-utils/graceful"
	"github.com/win-t/go-utils/http/defserver"
	"github.com/win-t/go-utils/mainutil/mainrun"
)

// Run http server.
func Run(setupFn func(*RunOption) error) {
	mainrun.Run(func() error {
		opt := RunOption{addr: ":8080"}

		if setupFn != nil {
			if err := setupFn(&opt); err != nil {
				return errors.Trace(err)
			}
		}

		s, err := defserver.New(opt.addr, opt.handler, opt.opts...)
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
	})
}

func wait(errCh chan error, s *http.Server) error {
	select {
	case err := <-errCh:
		return errors.Trace(err)
	case <-graceful.Context().Done():
	}

	timeout := 4 * time.Second
	if s.ReadTimeout > timeout {
		timeout = s.ReadTimeout
	}
	if s.ReadHeaderTimeout > timeout {
		timeout = s.ReadHeaderTimeout
	}
	if s.WriteTimeout > timeout {
		timeout = s.WriteTimeout
	}
	timeout += 1 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		return errors.Errorf(
			"cannot gracefuly shutdown server in %s",
			timeout.Truncate(time.Second).String(),
		)
	}

	return nil
}
