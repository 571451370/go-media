package sdl

import (
	"io"
	"os"
)

type RW interface {
	io.ReadWriteCloser
	io.Seeker
	Stat() os.FileInfo
}

type RWOps struct {
	RW
}
