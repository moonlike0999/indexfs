package internal

import "github.com/caddyserver/caddy/v2"

func ModGen[T any, PT interface {
	caddy.Module
	*T
}](name caddy.ModuleID) caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "caddy.fs." + name,
		New: func() caddy.Module {
			var t T
			var pt PT = &t
			return pt
		},
	}
}
