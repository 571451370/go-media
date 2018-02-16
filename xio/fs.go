package xio

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type FS interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	Stat(name string) (os.FileInfo, error)
}

type File interface {
	io.ReadWriteSeeker
	io.Closer
	Stat() (os.FileInfo, error)
	Readdir(int) ([]os.FileInfo, error)
}

type Stat struct {
	Path      string
	Length    int64
	Type      os.FileMode
	Mod       time.Time
	Interface interface{}
}

func (s Stat) Name() string       { return filepath.Base(s.Path) }
func (s Stat) Size() int64        { return s.Length }
func (s Stat) Mode() os.FileMode  { return s.Type }
func (s Stat) IsDir() bool        { return s.Type&os.ModeDir != 0 }
func (s Stat) ModTime() time.Time { return s.Mod }
func (s Stat) Sys() interface{}   { return s.Interface }

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

func (fs *SFS) Create(name string) (File, error) {
	name = filepath.Join(fs.Root, name)
	f, err := os.Create(name)
	return File(f), err
}

func (fs *SFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	name = filepath.Join(fs.Root, name)
	f, err := os.OpenFile(name, flag, perm)
	return File(f), err
}

func (fs *SFS) Remove(name string) error {
	name = filepath.Join(fs.Root, name)
	return os.Remove(name)
}

func (fs *SFS) Stat(name string) (os.FileInfo, error) {
	name = filepath.Join(fs.Root, name)
	return os.Stat(name)
}

func (fs *SFS) Rename(oldpath, newpath string) error {
	oldpath = filepath.Join(fs.Root, oldpath)
	newpath = filepath.Join(fs.Root, newpath)
	return os.Rename(oldpath, newpath)
}

func ReadFile(fs FS, name string) ([]byte, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func WriteFile(fs FS, name string, data []byte, perm os.FileMode) error {
	f, err := fs.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if xerr := f.Close(); err == nil {
		err = xerr
	}
	return err
}

func ReadAllDir(fs FS, name string) ([]os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fis, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	sort.Slice(fis, func(i, j int) bool {
		return fis[i].Name() < fis[j].Name()
	})

	return fis, nil
}
