package httprun

import (
	"context"
	"net/http"
	"time"

	"github.com/payfazz/go-errors/v2"
	"github.com/win-t/go-utils/graceful/gracefulhttp"
	"github.com/win-t/go-utils/http/defserver"
)

func Run(addr string, handler http.HandlerFunc) error {
	s, err := defserver.New(addr, handler)
	if err != nil {
		return errors.Trace(err)
	}

	gracefulErr := gracefulhttp.Shutdown(s, 5*time.Second)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Trace(err)
	}
	if err := gracefulErr(); err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return errors.New("cannot gracefuly shutdown server in 5 second")
		default:
			return errors.Trace(err)
		}
	}

	return nil
}
