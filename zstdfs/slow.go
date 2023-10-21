package zstdfs

import (
	"errors"
	"fmt"
	"github.com/klauspost/compress/zstd"
	"io"
	"io/fs"
	"time"
)

var (
	_ fs.File     = (*Slow)(nil)
	_ fs.DirEntry = (*Slow)(nil)
	_ fs.FileInfo = (*Slow)(nil)
	_ io.Seeker   = (*Slow)(nil)
)

type Slow struct {
	f   fs.File
	off int64

	open    func() (fs.File, error)
	dec     *zstd.Decoder
	name    string
	mode    fs.FileMode
	size    int64
	modTime time.Time
	sys     any
}

func (f *Slow) Read(b []byte) (int, error) {
	if err := f.init(); err != nil {
		return 0, err
	}

	n, err := f.dec.Read(b)
	f.off += int64(n)
	return n, err
}

func (f *Slow) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekCurrent:
		offset = f.off + offset
	case io.SeekEnd:
		offset = f.size + offset
	case io.SeekStart:
	default:
		return 0, errors.New("invalid whence")
	}

	if offset < 0 || offset > f.size {
		return 0, fmt.Errorf("offset %d is out of range", offset)
	} else if offset == f.off {
		return f.off, nil
	} else if offset > f.off {
		if _, err := io.CopyN(io.Discard, f, offset-f.off); err != nil {
			return 0, err
		}
		return f.off, nil
	}

	if err := f.reset(); err != nil {
		return 0, err
	}
	return f.Seek(offset, io.SeekStart)
}

func (f *Slow) init() error {
	if f.f == nil {
		f.off = 0
		fi, err := f.open()
		if err != nil {
			return err
		}

		f.f = fi
		if err := f.dec.Reset(f); err != nil {
			f.f = nil
			return errors.Join(err, fi.Close())
		}
	}
	return nil
}

func (f *Slow) reset() error {
	err := f.f.Close()
	f.f = nil
	if err != nil {
		return err
	}
	return f.init()
}

func (f *Slow) Close() error {
	f.dec.Close()
	if f.f == nil {
		return nil
	}
	return f.f.Close()
}

func (f *Slow) Size() int64                { return f.size }
func (f *Slow) Mode() fs.FileMode          { return f.mode }
func (f *Slow) ModTime() time.Time         { return f.modTime }
func (f *Slow) Sys() any                   { return f.sys }
func (f *Slow) Name() string               { return f.name }
func (f *Slow) IsDir() bool                { return false }
func (f *Slow) Type() fs.FileMode          { return f.mode.Type() }
func (f *Slow) Info() (fs.FileInfo, error) { return f, nil }
func (f *Slow) Stat() (fs.FileInfo, error) { return f, nil }
