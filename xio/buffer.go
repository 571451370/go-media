package xio

import "io"

type Buffer struct {
	buf []byte
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{
		buf: buf,
	}
}

func (b *Buffer) ReadAt(p []byte, off int64) (n int, err error) {
	if off >= int64(len(b.buf)) {
		return 0, io.EOF
	}
	n = copy(p, b.buf[off:])
	return
}

func (b *Buffer) WriteAt(p []byte, off int64) (n int, err error) {
	for i := range p {
		if off+int64(i) >= int64(len(b.buf)) {
			b.buf = append(b.buf, p[i])
		} else {
			b.buf[off+int64(i)] = p[i]
		}
	}
	return len(p), nil
}

func (b *Buffer) Buffer() []byte {
	return b.buf
}
