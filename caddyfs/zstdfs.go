package caddyfs

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/moonlike0999/indexfs/caddyfs/internal"
	"github.com/moonlike0999/indexfs/zstdfs"
	"io/fs"
)

var _ interface {
	caddy.Module
	fs.FS
	caddy.Provisioner
} = (*IndexFS)(nil)

func init() { caddy.RegisterModule(new(ZSTDFS)) }

type ZSTDFS struct {
	internal.BaseFS
}

func (zfs *ZSTDFS) CaddyModule() caddy.ModuleInfo {
	return internal.ModGen[ZSTDFS]("zstdfs")
}

func (zfs *ZSTDFS) Provision(ctx caddy.Context) error {
	if err := zfs.BaseFS.Provision(ctx); err != nil {
		return err
	}
	zfs.FS = zstdfs.New(zfs.FS)
	return nil
}
