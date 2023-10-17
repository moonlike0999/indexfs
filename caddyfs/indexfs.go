package caddyfs

import (
	"errors"
	"github.com/caddyserver/caddy/v2"
	"github.com/moonlike0999/indexfs/caddyfs/internal"
	"github.com/moonlike0999/indexfs/indexfs"
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
	allowed := make([]string, len(ifs.AllowedExtensions))
	replacer := caddy.NewReplacer()
	for i, s := range ifs.AllowedExtensions {
		allowed[i] = replacer.ReplaceAll(s, "")
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
