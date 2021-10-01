//go:build linux || darwin
// +build linux darwin

package mainhttp

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/payfazz/go-errors/v2"

	"github.com/win-t/go-utils/http/defserver"
	"github.com/win-t/go-utils/mainutil/mainrun"
)

// Run http server on unix socket.
func RunUnixSocket(setupFn func(*RunOption) error) {
	mainrun.Run(func() error {
		opt := RunOption{addr: "./http.sock"}

		if setupFn != nil {
			if err := setupFn(&opt); err != nil {
				return errors.Trace(err)
			}
		}

		var err error
		opt.addr, err = filepath.Abs(opt.addr)
		if err != nil {
			return errors.Trace(err)
		}

		listener, err := net.Listen("unix", opt.addr)
		if err != nil {
			return errors.Trace(err)
		}
		defer os.RemoveAll(opt.addr)

		s, err := defserver.New("", opt.handler, opt.opts...)
		if err != nil {
			return errors.Trace(err)
		}

		errCh := make(chan error, 1)
		go func() {
			fmt.Printf("Running http on unix socket %s\n", opt.addr)
			if err := s.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		return wait(errCh, s)
	})
}
