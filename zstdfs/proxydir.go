package zstdfs

import (
	"io"
	"io/fs"
	"strings"
	"syscall"
)

var (
	_ fs.ReadDirFile = (*proxyDir)(nil)
	_ io.Seeker      = (*proxyDir)(nil)
)

type proxyDir struct {
	fs.ReadDirFile
	base fs.FS
}

func (p *proxyDir) Seek(int64, int) (int64, error) { return 0, syscall.EISDIR }

func (p *proxyDir) ReadDir(count int) ([]fs.DirEntry, error) {
	entries, err := p.ReadDirFile.ReadDir(count)
	for i, entry := range entries {
		entries[i] = &proxyEntry{DirEntry: entry, base: p.base}
	}
	return entries, err
}

type proxyEntry struct {
	fs.DirEntry
	base fs.FS
}

func (p *proxyEntry) Name() string { return strings.TrimSuffix(p.DirEntry.Name(), ".zst") }

func (p *proxyEntry) Info() (fs.FileInfo, error) {
	stat, err := p.DirEntry.Info()
	if err != nil {
		return nil, err
	}
	proxy := &proxyInfo{FileInfo: stat, name: p.Name()}

	stat, err = fs.Stat(p.base, p.Name())
	if err != nil {
		return nil, err
	}
	proxy.size = stat.Size()

	return proxy, nil
}

type proxyInfo struct {
	fs.FileInfo
	name string
	size int64
}

func (p *proxyInfo) Name() string { return p.name }
func (p *proxyInfo) Size() int64  { return p.size }
