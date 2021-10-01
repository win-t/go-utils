package httprun

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/payfazz/go-errors/v2"
)

func RunUnixSocket(ctx context.Context, s *http.Server) error {
	if s.Addr == "" {
		s.Addr = "./http.sock"
	}
	socket, err := filepath.Abs(s.Addr)
	if err != nil {
		return errors.Trace(err)
	}
	s.Addr = ""

	listener, err := net.Listen("unix", socket)
	if err != nil {
		return errors.Trace(err)
	}
	defer os.RemoveAll(socket)

	errCh := make(chan error, 1)
	go func() { errCh <- ignoreErrServerClosed(s.Serve(listener)) }()

	return wait(ctx, errCh, s)
}
