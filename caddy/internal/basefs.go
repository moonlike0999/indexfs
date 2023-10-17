package internal

import (
	"encoding/json"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"io/fs"
)

var _ interface {
	caddy.Provisioner
	fs.FS
} = (*BaseFS)(nil)

type BaseFS struct {
	FS            fs.FS           `json:"-"`
	FileSystemRaw json.RawMessage `json:"file_system,omitempty" caddy:"namespace=caddy.fs inline_key=backend"`
}

func (bfs *BaseFS) Provision(ctx caddy.Context) error {
	if len(bfs.FileSystemRaw) > 0 {
		mod, err := ctx.LoadModule(bfs, "FileSystemRaw")
		if err != nil {
			return fmt.Errorf("loading file system module: %v", err)
		}
		bfs.FS = mod.(fs.FS)
	} else {
		return fmt.Errorf("no fs specified")
	}
	return nil
}

func (bfs *BaseFS) Open(name string) (fs.File, error) {
	return bfs.FS.Open(name)
}
