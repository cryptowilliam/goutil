package gchacha20

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"log"
	"os"
	"testing"
)

func TestChaCha20RWC_Read_Write(t *testing.T) {
	plain := "0123456789"
	filepath := "chacha20_test.bin"
	passphrase := "this-is-a-passphrase"

	makerW, err := NewChaCha20Maker(passphrase)
	gtest.Assert(t, err)
	tempFile, err := os.CreateTemp("", filepath)
	gtest.Assert(t, gerrors.Wrap(err, "create temp error"))
	defer tempFile.Close()
	chachaW, err := makerW.Make(tempFile, false, nil)
	gtest.Assert(t, gerrors.Wrap(err, "make error"))
	n, err := chachaW.Write([]byte(plain))
	gtest.Assert(t, gerrors.Wrap(err, "write error"))
	log.Println("write size", n)


	_, err = tempFile.Seek(0, 0)
	gtest.Assert(t, gerrors.Wrap(err, "seek error"))
	makerR, err := NewChaCha20Maker(passphrase)
	gtest.Assert(t, err)
	gtest.Assert(t, gerrors.Wrap(err, "open temp error"))
	chachaR, err := makerR.Make(tempFile, true, nil)
	gtest.Assert(t, gerrors.Wrap(err, "m2 make error"))
	var readBuf = make([]byte, 1024)
	n, err = chachaR.Read(readBuf)
	gtest.Assert(t, gerrors.Wrap(err, "read error"))
	log.Println("read size", n)
	log.Println("read content", string(readBuf[:n]))
	if string(readBuf[:n]) != plain {
		gtest.Assert(t, gerrors.New("chacha20 read write error"))
	}
}
