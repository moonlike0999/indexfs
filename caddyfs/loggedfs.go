package caddyfs

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/moonlike0999/indexfs/caddyfs/internal"
	"github.com/moonlike0999/indexfs/loggedfs"
	"go.mrchanchal.com/zaphandler"
	"io/fs"
	"log/slog"
)

var _ interface {
	caddy.Module
	caddy.Provisioner
	fs.FS
} = (*LoggedFS)(nil)

func init() { caddy.RegisterModule(new(LoggedFS)) }

type LoggedFS struct {
	internal.BaseFS
}

func (ifs *LoggedFS) CaddyModule() caddy.ModuleInfo {
	return internal.ModGen[LoggedFS]("loggedfs")
}

func (ifs *LoggedFS) Provision(ctx caddy.Context) error {
	if err := ifs.BaseFS.Provision(ctx); err != nil {
		return err
	}
	ifs.FS = loggedfs.Logged(slog.New(zaphandler.New(ctx.Logger())), ifs.FS)
	return nil
}

func (ifs *LoggedFS) Open(name string) (fs.File, error) {
	return ifs.FS.Open(name)
}
