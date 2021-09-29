//go:build linux || darwin
// +build linux darwin

package httprun

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/payfazz/go-errors/v2"
	"github.com/win-t/go-utils/graceful/gracefulhttp"
	"github.com/win-t/go-utils/http/defserver"
)

func RunOnUnixSocket(socket string, handler http.HandlerFunc) error {
	socket, err := filepath.Abs(socket)
	if err != nil {
		return errors.Trace(err)
	}

	listener, err := net.Listen("unix", socket)
	if err != nil {
		return errors.Trace(err)
	}
	defer os.RemoveAll(socket)

	s, err := defserver.New("", handler)
	if err != nil {
		return errors.Trace(err)
	}

	gracefulErr := gracefulhttp.Shutdown(s, 5*time.Second)
	if err := s.Serve(listener); err != nil && err != http.ErrServerClosed {
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
