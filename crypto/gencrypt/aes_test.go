package gencrypt

import (
	"bytes"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"github.com/cryptowilliam/goutil/container/grand"
	"github.com/sekrat/aescrypter"
	"testing"
)

func TestAesGcm256_Encrypt(t *testing.T) {
	for i := 0; i < 100; i++ {
		plain := grand.RandomString(grand.RandomInt(0, 1000))
		secret := grand.RandomString(grand.RandomInt(0, 1000))
		cipher, err := NewAesGcm256().Encrypt([]byte(plain), []byte(secret), true)
		gtest.Assert(t, err)
		decryptPlain, err := aescrypter.New().Decrypt(secret, cipher)
		gtest.Assert(t, err)
		if !bytes.Equal([]byte(plain), decryptPlain) {
			gtest.PrintlnExit(t, "origin %v, decrypt %v", plain, decryptPlain)
		}
	}
}

func TestAesGcm256_Decrypt(t *testing.T) {
	for i := 0; i < 100; i++ {
		plain := grand.RandomString(grand.RandomInt(0, 1000))
		secret := grand.RandomString(grand.RandomInt(0, 1000))
		cipher, err := aescrypter.New().Encrypt(secret, []byte(plain))
		gtest.Assert(t, err)
		decryptPlain, err := NewAesGcm256().Decrypt(cipher, []byte(secret), true)
		gtest.Assert(t, err)
		if !bytes.Equal([]byte(plain), decryptPlain) {
			gtest.PrintlnExit(t, "origin %v, decrypt %v", plain, decryptPlain)
		}
	}
}
