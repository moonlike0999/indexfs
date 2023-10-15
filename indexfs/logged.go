package indexfs

import (
	"io/fs"
	"log/slog"
)

var (
	_ fs.FS = (*_Logged)(nil)
)

type _Logged struct {
	Logger *slog.Logger
	FS     fs.FS
}

func Logged(logger *slog.Logger, fsys fs.FS) fs.FS {
	return &_Logged{logger, fsys}
}

func (l *_Logged) Open(name string) (fs.File, error) {
	f, err := l.FS.Open(name)
	if err != nil {
		l.Logger.Error("failed to open file", "name", name, "err", err)
	} else {
		l.Logger.Info("opened file", "name", name)
	}
	return f, err
}
