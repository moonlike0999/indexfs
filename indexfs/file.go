package indexfs

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	_ fs.File     = (*File)(nil)
	_ fs.DirEntry = (*File)(nil)
	_ fs.FileInfo = (*File)(nil)
	_ io.Seeker   = (*File)(nil)
)

type File struct {
	_FS              fs.FS
	_File            fs.File
	_OpenErr         error
	_Opened, _Closed sync.Once
	Path             string
	Date             *Date
	FileSize         int64
}

func (f *File) _Open() error {
	f._Opened.Do(func() {
		f._File, f._OpenErr = f._FS.Open(f.Path)
		f._FS = nil
	})
	return f._OpenErr
}

func (f *File) Read(b []byte) (int, error) {
	if err := f._Open(); err != nil {
		return 0, err
	}
	return f._File.Read(b)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if err := f._Open(); err != nil {
		return 0, err
	}
	seeker, ok := f._File.(io.Seeker)
	if !ok {
		return 0, errors.New("does not implement seeker")
	}
	return seeker.Seek(offset, whence)
}

func (f *File) Close() (err error) {
	f._Closed.Do(func() {
		if err = f._Open(); err == nil {
			err = f._File.Close()
		}
	})
	return
}

func (f *File) Stat() (fs.FileInfo, error) { return f, nil }
func (f *File) Name() string               { return f.Date.String() + filepath.Ext(f.Path) }
func (f *File) Size() int64                { return f.FileSize }
func (f *File) Mode() fs.FileMode          { return os.ModePerm }
func (f *File) Type() fs.FileMode          { return f.Mode().Type() }
func (f *File) ModTime() time.Time         { return f.Date.Time() }
func (f *File) IsDir() bool                { return false }
func (f *File) Sys() any                   { return nil }
func (f *File) Info() (fs.FileInfo, error) { return f, nil }
