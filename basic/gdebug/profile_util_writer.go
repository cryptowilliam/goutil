package gdebug

import "io"

type writer struct {
	io.Writer
}

func (r *writer) Open(name string) (io.WriteCloser, error) {
	return r, nil
}

func (r *writer) Close() error {
	return nil
}

func newWriter(w io.Writer) *writer {
	return &writer{w}
}
