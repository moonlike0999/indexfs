package zstdfs

import (
	"errors"
	"github.com/klauspost/compress/zstd"
	"gitlab.com/rackn/seekable-zstd"
	"io"
	"io/fs"
	"strings"
	"time"
)

var (
	_ fs.File     = (*Fast)(nil)
	_ fs.DirEntry = (*Fast)(nil)
	_ fs.FileInfo = (*Fast)(nil)
	_ io.Seeker   = (*Fast)(nil)
)

type Fast struct {
	f    fs.File
	stat fs.FileInfo
	dec  *zstd.Decoder
	sdec *seekable.Decoder
	rs   io.ReadSeeker
}

func (f *Fast) Seek(offset int64, whence int) (int64, error) { return f.rs.Seek(offset, whence) }
func (f *Fast) Size() int64                                  { return f.sdec.Size() }
func (f *Fast) Mode() fs.FileMode                            { return f.stat.Mode() }
func (f *Fast) ModTime() time.Time                           { return f.stat.ModTime() }
func (f *Fast) Sys() any                                     { return f.stat.Sys() }
func (f *Fast) Name() string                                 { return strings.TrimSuffix(f.stat.Name(), ".zst") }
func (f *Fast) IsDir() bool                                  { return false }
func (f *Fast) Type() fs.FileMode                            { return f.Mode().Type() }
func (f *Fast) Info() (fs.FileInfo, error)                   { return f, nil }
func (f *Fast) Stat() (fs.FileInfo, error)                   { return f, nil }
func (f *Fast) Read(b []byte) (int, error)                   { return f.rs.Read(b) }
func (f *Fast) Close() error                                 { f.dec.Close(); return errors.Join(f.sdec.Close(), f.f.Close()) }
