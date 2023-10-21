package caddyfs

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/moonlike0999/indexfs/caddyfs/internal"
	"github.com/moonlike0999/indexfs/indexfs"
	"io/fs"
)

var _ interface {
	caddy.Module
	fs.FS
	caddy.Provisioner
} = (*IndexFS)(nil)

func init() { caddy.RegisterModule(new(IndexFS)) }

type IndexFS struct {
	internal.BaseFS
}

func (ifs *IndexFS) CaddyModule() caddy.ModuleInfo {
	return internal.ModGen[IndexFS]("indexfs")
}

func (ifs *IndexFS) Provision(ctx caddy.Context) error {
	if err := ifs.BaseFS.Provision(ctx); err != nil {
		return err
	}
	ifs.FS = indexfs.New(ifs.FS)
	return nil
}
