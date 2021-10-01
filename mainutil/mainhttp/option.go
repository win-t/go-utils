package mainhttp

import (
	"net/http"

	"github.com/win-t/go-utils/http/defserver"
)

type RunOption struct {
	addr    string
	handler http.HandlerFunc
	opts    []defserver.Option
}

func (r *RunOption) SetAddr(addr string) {
	r.addr = addr
}

func (r *RunOption) SetHandler(handler http.HandlerFunc) {
	r.handler = handler
}

func (r *RunOption) SetServerOptions(opts ...defserver.Option) {
	r.opts = opts
}
