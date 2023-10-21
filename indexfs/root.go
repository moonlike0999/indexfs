package indexfs

import (
	"errors"
	"golang.org/x/exp/maps"
	"io"
	"io/fs"
	"os"
	"regexp"
	"slices"
	"syscall"
	"time"
)

var (
	_ fs.FileInfo    = (*_Root)(nil)
	_ fs.DirEntry    = (*_Root)(nil)
	_ fs.File        = (*_Root)(nil)
	_ fs.ReadDirFile = (*_Root)(nil)
	_ io.Seeker      = (*_Root)(nil)
)

type _Root FS

func _NewRoot(fsys *FS) fs.File { return (*_Root)(fsys) }

func (r *_Root) ReadDir(count int) ([]fs.DirEntry, error) {
	files, err := _IndexFiles(r._Base, r._FileRegex)
	return _LimitFiles(count, files), err
}

func _IndexFiles(base fs.FS, fileRegex *regexp.Regexp) ([]fs.DirEntry, error) {
	files := map[Date]fs.DirEntry{}
	err := _WalkFiles(base, fileRegex, func(f *File) error { files[*f.Date] = f; return nil })
	return _SortFileInfo(files), err
}

func _WalkFiles(base fs.FS, fileRegex *regexp.Regexp, handleFile func(*File) error) error {
	var errs []error
	errs = append(errs, fs.WalkDir(base, ".", _WalkFunc(base, fileRegex, &errs, handleFile)))
	return errors.Join(errs...)
}

func _LimitFiles(count int, files []fs.DirEntry) []fs.DirEntry {
	if count >= 0 && count < len(files) {
		return files[:count]
	}
	return files
}

func _WalkFunc(fsys fs.FS, fileRegex *regexp.Regexp, errs *[]error, handleFile func(*File) error) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			*errs = append(*errs, err)
		} else if !d.IsDir() && d.Type().IsRegular() && fileRegex.MatchString(path) {
			f, err := _GetFileInfo(fsys, path, d)
			if err != nil {
				*errs = append(*errs, err)
			}
			return handleFile(f)
		}
		return nil
	}
}

func _SortFileInfo(files map[Date]fs.DirEntry) []fs.DirEntry {
	a := maps.Values(files)
	slices.SortFunc(a, func(a, b fs.DirEntry) int {
		af, bf := a.(*File), b.(*File)
		if af.Date.Year == bf.Date.Year {
			return int(af.Date.Month) - int(bf.Date.Month)
		}
		return int(af.Date.Year) - int(bf.Date.Year)
	})
	return a
}

func _GetFileInfo(fsys fs.FS, path string, d fs.DirEntry) (_ *File, finalErr error) {
	f := File{_FS: fsys, Path: path, Date: new(Date)}

	if stat, err := d.Info(); err != nil {
		f.FileSize = -1
		finalErr = errors.Join(finalErr, err)
	} else {
		f.FileSize = stat.Size()
	}

	if err := f.Date.UnmarshalText([]byte(path)); err != nil {
		finalErr = errors.Join(finalErr, err)
		f.Date = new(Date)
	}
	return &f, finalErr
}

func (r *_Root) ModTime() time.Time {
	files, _ := r.ReadDir(-1)
	if len(files) == 0 {
		return time.Time{}
	}
	return files[len(files)-1].(*File).Date.Time()
}

func (r *_Root) Read([]byte) (n int, err error) { return 0, syscall.EISDIR }
func (r *_Root) Close() error                   { return nil }
func (r *_Root) Seek(int64, int) (int64, error) { return 0, syscall.EISDIR }
func (r *_Root) Stat() (fs.FileInfo, error)     { return r, nil }
func (r *_Root) Name() string                   { return "/" }
func (r *_Root) Size() int64                    { return 0 }
func (r *_Root) Mode() fs.FileMode              { return os.ModeDir | os.ModePerm }
func (r *_Root) IsDir() bool                    { return true }
func (r *_Root) Sys() any                       { return nil }
func (r *_Root) Type() fs.FileMode              { return r.Mode().Type() }
func (r *_Root) Info() (fs.FileInfo, error)     { return r, nil }
