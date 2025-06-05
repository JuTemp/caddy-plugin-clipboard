package caddy_plugin

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

// https://caddyserver.com/docs/extending-caddy#complete-example

type Clipboard struct {
	logger *zap.Logger
}

func init() {
	caddy.RegisterModule(Clipboard{})
	httpcaddyfile.RegisterHandlerDirective("clipboard", parseCaddyfile)
}

func (Clipboard) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.clipboard",
		New: func() caddy.Module { return new(Clipboard) },
	}
}

func (m *Clipboard) Provision(ctx caddy.Context) error {
	m.logger = ctx.Logger(m)
	return nil
}

func (m *Clipboard) Validate() error {
	return nil
}

func (m Clipboard) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return ServeHttp(w, r, next, m)
}

func (m *Clipboard) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Clipboard
	return m, nil
}

var (
	_ caddy.Provisioner           = (*Clipboard)(nil)
	_ caddy.Validator             = (*Clipboard)(nil)
	_ caddyhttp.MiddlewareHandler = (*Clipboard)(nil)
	_ caddyfile.Unmarshaler       = (*Clipboard)(nil)
)
