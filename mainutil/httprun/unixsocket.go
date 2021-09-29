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
	"github.com/win-t/go-utils/graceful"
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

	errCh := make(chan error, 1)
	go func() {
		if err := s.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return errors.Trace(err)
	case <-graceful.Context().Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(shutdownCtx); err != nil {
		return errors.New("cannot gracefuly shutdown server in 5 second")
	}

	return nil
}
