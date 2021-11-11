package gio

import (
	"io"
)

type (
	Flusher interface {
		Flush() error
	}

	WriteFlusher interface {
		io.Writer
		Flusher
	}
)
