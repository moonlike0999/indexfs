package indexfs

import (
	"errors"
	"github.com/caddyserver/caddy/v2"
	"github.com/moonlike0999/caddyfs/indexfs"
	"go.mrchanchal.com/zaphandler"
	"io/fs"
	"log/slog"
	"os"
)

var _ caddy.Module = (*IndexFS)(nil)

func init() { caddy.RegisterModule(new(IndexFS)) }

type IndexFS struct {
	fs.FS             `json:"-"`
	DataDir           string   `json:"dir"`
	AllowedExtensions []string `json:"extensions"`
}

func (ifs *IndexFS) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.fs.indexfs",
		New: func() caddy.Module { return new(IndexFS) },
	}
}

func (ifs *IndexFS) Provision(ctx caddy.Context) error {
	ifs.FS = indexfs.New(os.DirFS(ifs.DataDir), ifs.AllowedExtensions...)
	ifs.FS = indexfs.Logged(slog.New(zaphandler.New(ctx.Logger())), ifs.FS)
	return nil
}

func (ifs *IndexFS) Validate() error {
	if len(ifs.AllowedExtensions) == 0 {
		return errors.New("no allowed extensions")
	} else if _, err := os.Stat(ifs.DataDir); err != nil {
		return err
	}
	return nil
}
