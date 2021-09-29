// Package defmiddleware.
package defmiddleware

import (
	"net/http"

	"github.com/payfazz/go-middleware/panicreporter"
	"github.com/payfazz/go-middleware/reqlogger"
)

// Get return middleware that commonly used.
func Get() []func(http.HandlerFunc) http.HandlerFunc {
	return []func(http.HandlerFunc) http.HandlerFunc{
		panicreporter.New(nil),
		reqlogger.New(nil),
	}
}
