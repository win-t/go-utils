// Package defserver.
package defserver

import (
	"net/http"
	"time"

	"github.com/payfazz/go-middleware"

	"github.com/win-t/go-utils/http/defmiddleware"
)

type Option func(*http.Server) error

// New return http.Server that commonly used.
func New(addr string, handler http.HandlerFunc, opts ...Option) (*http.Server, error) {
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
