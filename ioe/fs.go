package ioe

import (
	"io"
	"os"
	"path/filepath"
)

type FS interface {
	Open(name string) (File, error)
}

type File interface {
	io.ReadWriteSeeker
	io.Closer
	Stat() (os.FileInfo, error)
}

type SFS struct {
	Root string
}

func (fs *SFS) Chdir(dir string) error {
	fs.Root = filepath.Join(fs.Root, dir)
	return nil
}

func (fs *SFS) Open(name string) (File, error) {
	name = filepath.Join(fs.Root, name)
	f, err := os.Open(name)
	return File(f), err
}
