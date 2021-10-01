package httprun

import (
	"context"
	"net/http"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = errors.Errorf("http shutdown timeout: %w", err)
		}
		return err
	}

	return nil
}
