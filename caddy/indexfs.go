package caddyfs

import (
	"errors"
	"github.com/caddyserver/caddy/v2"
	"github.com/moonlike0999/caddyfs/indexfs"
	"github.com/moonlike0999/indexfs/caddy/internal"
	"io/fs"
)

var _ interface {
	caddy.Module
	fs.FS
	caddy.Provisioner
	caddy.Validator
} = (*IndexFS)(nil)

func init() { caddy.RegisterModule(new(IndexFS)) }

type IndexFS struct {
	internal.BaseFS
	AllowedExtensions []string `json:"extensions"`
}

func (ifs *IndexFS) CaddyModule() caddy.ModuleInfo {
	return internal.ModGen[IndexFS]("indexfs")
}

func (ifs *IndexFS) Provision(ctx caddy.Context) error {
	if err := ifs.BaseFS.Provision(ctx); err != nil {
		return err
	}
	ifs.FS = indexfs.New(ifs.FS, ifs.AllowedExtensions...)
	return nil
}

func (ifs *IndexFS) Validate() error {
	if len(ifs.AllowedExtensions) == 0 {
		return errors.New("no allowed extensions")
	}
	return nil
}
