package alidns

import (
	caddy "github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/alidns"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *alidns.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.alidns",
		New: func() caddy.Module { return &Provider{new(alidns.Provider)} },
	}
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// alidns {
//     access_key_id "<access_key_id>"
//     access_key_secret "<access_key_secret>"
// }
//
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	repl := caddy.NewReplacer()
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "access_key_id":
				if d.NextArg() {
					p.Provider.AccKeyID = repl.ReplaceAll(d.Val(), "")
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "access_key_secret":
				if d.NextArg() {
					p.Provider.AccKeySecret = repl.ReplaceAll(d.Val(), "")
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.AccKeyID == "" || p.AccKeySecret == "" {
		return d.Err("AccessKeyID or AccessKeySecret is empty")
	}
	return nil
}

// Interface guard
var _ caddyfile.Unmarshaler = (*Provider)(nil)
