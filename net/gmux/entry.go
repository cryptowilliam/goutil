package gmux

import (
	"io"
)

func NewClient(conn io.ReadWriteCloser, config *Config) (*Session, error) {
	if config == nil {
		config = DefaultConfig()
	}
	return Client(conn, config)
}

func NewServer(conn io.ReadWriteCloser, config *Config) (*Session, error) {
	if config == nil {
		config = DefaultConfig()
	}
	return Server(conn, config)
}