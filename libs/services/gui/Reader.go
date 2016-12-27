package gui

import (
	"io"
)

type Reader struct {
	io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{Reader: r}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return 0, nil
}
