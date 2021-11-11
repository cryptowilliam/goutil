package gmono256

import (
	"crypto/rand"
	"reflect"
	"sort"
	"testing"
)

const (
	_MB_ = 1024 * 1024
)

func TestRandAlphabet(t *testing.T) {
	abt := randAlphabet()
	t.Log(abt)
	b64sAbt, err := Base64ToAlphabet(abt)
	if err != nil {
		t.Error(err)
	}
	sort.Sort(b64sAbt)
	for i := 0; i < alphabetLen; i++ {
		if b64sAbt[i] != byte(i) {
			t.Error("Duplicate byte found, it is not allowed.")
		}
	}
}

func TestMono256Cipher(t *testing.T) {
	abt := randAlphabet()
	t.Log(abt)
	p, _ := Base64ToAlphabet(abt)
	cipher := NewMono256(p)
	// plaintext
	org := make([]byte, alphabetLen)
	for i := 0; i < alphabetLen; i++ {
		org[i] = byte(i)
	}
	// copy plaintext to tmp
	tmp := make([]byte, alphabetLen)
	copy(tmp, org)
	t.Log(tmp)
	// encode tmp
	cipher.Encrypt(tmp)
	t.Log(tmp)
	// decode tmp
	cipher.Decrypt(tmp)
	t.Log(tmp)
	if !reflect.DeepEqual(org, tmp) {
		t.Error("Data does not correspond after decoding.")
	}
}

func BenchmarkMono256Cipher_Encode(b *testing.B) {
	abt := randAlphabet()
	p, _ := Base64ToAlphabet(abt)
	cipher := NewMono256(p)
	bs := make([]byte, _MB_)
	b.ResetTimer()
	rand.Read(bs)
	cipher.Encrypt(bs)
}

func BenchmarkMono256Cipher_Decode(b *testing.B) {
	abt := randAlphabet()
	p, _ := Base64ToAlphabet(abt)
	cipher := NewMono256(p)
	bs := make([]byte, _MB_)
	b.ResetTimer()
	rand.Read(bs)
	cipher.Decrypt(bs)
}
