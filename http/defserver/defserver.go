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
func New(addr string, handler http.HandlerFunc, opts ...Option) (*http.Server, error) {
	if handler == nil {
		handler = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(501)
			w.Write([]byte("501 Not Implemented"))
		}
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
