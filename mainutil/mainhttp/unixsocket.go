//go:build linux || darwin
// +build linux darwin

package mainhttp

import (
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/payfazz/go-errors/v2"

	"github.com/win-t/go-utils/http/defserver"
	"github.com/win-t/go-utils/mainutil/mainrun"
)

// Run http server on unix socket from handler returned by getHandler
func RunUnixSocket(socket string, getHandler func() (http.HandlerFunc, error)) {
	mainrun.Run(func() error {
		handler, err := handlerOf(getHandler)
		if err != nil {
			return errors.Trace(err)
		}

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

		return wait(errCh, s)
	})
}
