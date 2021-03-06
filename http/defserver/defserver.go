// Package defserver.
package defserver

import (
	"net/http"
	"time"

	"github.com/payfazz/go-middleware"

	"github.com/win-t/go-utils/deftls"
	"github.com/win-t/go-utils/http/defmiddleware"
)

type Option func(*http.Server) error

// New return http.Server that commonly used.
func New(addr string, handler http.Handler, opts ...Option) (*http.Server, error) {
	if handler == nil {
		handler = http.DefaultServeMux
	}

	server := &http.Server{
		Addr: addr,
		Handler: middleware.C(
			defmiddleware.Get(),
			handler,
		),
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  30 * time.Second,
	}

	for _, o := range opts {
		if o != nil {
			if err := o(server); err != nil {
				return server, err
			}
		}
	}

	if _, ok := server.Handler.(noDefaultMiddleware); ok {
		server.Handler = handler
	}

	return server, nil
}

func WithTLS(opts ...deftls.Option) Option {
	return func(s *http.Server) error {
		config, err := deftls.Config(opts...)
		if err != nil {
			return err
		}
		if len(config.Certificates) == 0 {
			if err := deftls.UseCertSelfSigned()(config); err != nil {
				return err
			}
		}
		config.NextProtos = []string{"h2", "http/1.1"}
		s.TLSConfig = config
		return nil
	}
}

func WithNoTimeout() Option {
	return func(s *http.Server) error {
		s.ReadTimeout = 0
		s.WriteTimeout = 0
		s.IdleTimeout = 0
		return nil
	}
}

type noDefaultMiddleware struct{}

func (noDefaultMiddleware) ServeHTTP(http.ResponseWriter, *http.Request) {}

func WithNoDefaultMiddleware() Option {
	return func(s *http.Server) error {
		s.Handler = noDefaultMiddleware{}
		return nil
	}
}
