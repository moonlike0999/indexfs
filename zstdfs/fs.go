package zstdfs

import (
	"errors"
	"github.com/klauspost/compress/zstd"
	"gitlab.com/rackn/seekable-zstd"
	"io"
	"io/fs"
	"strings"
)

var (
	_ fs.FS = (*FS)(nil)
)

type FS struct {
	_Base         fs.FS
	slowSizeCache map[string]int64
}

func New(base fs.FS) *FS {
	return &FS{
		_Base:         base,
		slowSizeCache: make(map[string]int64),
	}
}

func (fsys *FS) Open(name string) (fs.File, error) {
	f, err := fsys._Base.Open(name)
	if e := err; e != nil {
		if f, err = fsys._Base.Open(name + ".zst"); err != nil {
			return nil, errors.Join(err, e)
		}
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, errors.Join(err, f.Close())
	} else if stat.IsDir() {
		return &proxyDir{ReadDirFile: f.(fs.ReadDirFile), base: fsys}, nil
	}

	dec, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1))
	if err != nil {
		return nil, errors.Join(err, f.Close())
	}

	sdec, sderr := seekable.NewDecoder(&readSeekerAt{r: f.(io.ReadSeeker)}, stat.Size(), dec)
	if sderr == nil {
		return &Fast{
			f:    f,
			stat: stat,
			dec:  dec,
			sdec: sdec,
			rs:   sdec.ReadSeeker(),
		}, nil
	}

	if _, err = f.(io.Seeker).Seek(0, io.SeekStart); err != nil {
		dec.Close()
		return nil, err
	}
	if err := dec.Reset(f); err != nil {
		dec.Close()
		return nil, errors.Join(sderr, err)
	}

	var size int64
	if ct, ok := fsys.slowSizeCache[name]; ok {
		size = ct
	} else {
		var ct countDiscarded
		if _, err = io.Copy(&ct, dec); err != nil {
			dec.Close()
			return nil, err
		}
		size = int64(ct)
	}
	fsys.slowSizeCache[name] = size

	if err := f.Close(); err != nil {
		dec.Close()
		return nil, err
	}

	return &Slow{
		open:    func() (fs.File, error) { return fsys._Base.Open(name) },
		dec:     dec,
		name:    strings.TrimSuffix(stat.Name(), ".zst"),
		mode:    stat.Mode(),
		size:    size,
		modTime: stat.ModTime(),
		sys:     stat.Sys(),
	}, nil
}

type countDiscarded int64

func (cd *countDiscarded) Write(b []byte) (int, error) {
	*cd += countDiscarded(len(b))
	return len(b), nil
}

type readSeekerAt struct {
	r io.ReadSeeker
}

func (r *readSeekerAt) ReadAt(p []byte, off int64) (n int, err error) {
	if _, err = r.r.Seek(off, io.SeekStart); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}
