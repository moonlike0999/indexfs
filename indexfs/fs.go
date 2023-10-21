package indexfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"slices"
)

var (
	_ fs.FS = (*FS)(nil)
)

type FS struct {
	_FileRegex *regexp.Regexp
	_Base      fs.FS
}

func New(base fs.FS) *FS {
	return &FS{
		_FileRegex: _MakeFileRegex(),
		_Base:      base,
	}
}

func _MakeFileRegex() *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`(?i)^.*%s.*$`, DateRegexString))
}

func (fsys *FS) Open(name string) (fs.File, error) {
	if slices.Contains([]string{"/", ".", ""}, name) {
		return _NewRoot(fsys), nil
	}

	var in Date
	if err := in.UnmarshalText([]byte(name)); err != nil {
		return nil, errors.Join(err, os.ErrNotExist)
	}

	return _OpenFile(fsys._Base, fsys._FileRegex, in)
}

func _OpenFile(base fs.FS, fileRegex *regexp.Regexp, date Date) (fs.File, error) {
	found, errC := make(chan fs.File), make(chan error)
	go func() {
		defer close(found)
		if err := _WalkFiles(base, fileRegex, func(file *File) error {
			if *file.Date == date {
				found <- file
				return fs.SkipAll
			}
			return nil
		}); err != nil {
			errC <- err
		}
	}()

	select {
	case err := <-errC:
		return nil, err
	case f, ok := <-found:
		if ok {
			return f, nil
		}
		return nil, os.ErrNotExist
	}
}
