// vwriter (void writer) is a package with a Writer implementing a write
// that returns nothing.
// Opposed to io.Writer which returns an error and number bytes written.
// Any errors should be handled.
package vwrite

import (
	"io"
)

type Writer interface {
	Write(p Handler)        // Write directs the io.Writer to use the handler provided.
	EWrite(p Handler) error // EWrite returns an error if any, useful if no error expected but to simplify return.
}

type writer struct {
	io.Writer
}

func New(w io.Writer) Writer {
	return &writer{Writer: w}
}

func (wr *writer) Write(p Handler) {
	p.Write(wr.Writer)
}

func (wr *writer) EWrite(p Handler) error {
	return p.EWrite(wr.Writer)
}

type Handler interface {
	Write(w io.Writer)
	EWrite(w io.Writer) error
}

type HandlerFunc func(w io.Writer)

func (p HandlerFunc) Write(w io.Writer) {
	p(w)
}

func (p HandlerFunc) EWrite(w io.Writer) error {
	p(w)
	return nil
}
