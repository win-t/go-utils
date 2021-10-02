package httprun

import (
	"context"
	"net/http"

	"github.com/payfazz/go-errors/v2"
)

func ignoreErrServerClosed(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func wait(ctx context.Context, errCh chan error, s *http.Server) error {
	select {
	case err := <-errCh:
		return errors.Trace(err)
	case <-ctx.Done():
	}

	if err := s.Shutdown(context.Background()); err != nil {
		return errors.Trace(err)
	}

	return nil
}
