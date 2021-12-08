package gchacha20

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"log"
	"os"
	"testing"
)

func TestChaCha20RWC_Read(t *testing.T) {
	m, err := NewChaCha20Maker("this-is-a-passphrase")
	gtest.Assert(t, err)

	plain := "0123456789"
	filepath := "chacha20_test.bin"

	tempFile, err := os.CreateTemp("", filepath)
	gtest.Assert(t, gerrors.Wrap(err, "create temp error"))
	defer tempFile.Close()
	chacha, err := m.Make(tempFile)
	gtest.Assert(t, gerrors.Wrap(err, "make error"))
	n, err := chacha.Write([]byte(plain))
	gtest.Assert(t, gerrors.Wrap(err, "write error"))
	log.Println("write size", n)

	_, err = tempFile.Seek(0, 0)
	gtest.Assert(t, gerrors.Wrap(err, "seek error"))
	var readBuf = make([]byte, 1024)
	n, err = chacha.Read(readBuf)
	gtest.Assert(t, gerrors.Wrap(err, "read error"))
	log.Println("read size", n)
	log.Println("read content", string(readBuf[:n]))
	if string(readBuf[:n]) != plain {
		gtest.Assert(t, gerrors.New("chacha20 read write error"))
	}
}
