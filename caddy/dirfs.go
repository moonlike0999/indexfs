package caddyfs

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/moonlike0999/indexfs/caddy/internal"
	"io/fs"
	"os"
)

var _ interface {
	caddy.Module
	fs.FS
	caddy.Provisioner
	caddy.Validator
} = (*DirFS)(nil)

func init() { caddy.RegisterModule(new(DirFS)) }

type DirFS struct {
	FS      fs.FS  `json:"-"`
	DataDir string `json:"dir"`
}

func (dfs *DirFS) CaddyModule() caddy.ModuleInfo {
	return internal.ModGen[DirFS]("dirfs")
}

func (dfs *DirFS) Provision(caddy.Context) error {
	dfs.FS = os.DirFS(dfs.DataDir)
	return nil
}

func (dfs *DirFS) Validate() error {
	if _, err := dfs.FS.(fs.StatFS).Stat("."); err != nil {
		return err
	}
	return nil
}

func (dfs *DirFS) Open(name string) (fs.File, error) {
	return dfs.FS.Open(name)
}
