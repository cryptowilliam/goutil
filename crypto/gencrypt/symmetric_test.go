package gencrypt

import (
	"bytes"
	"github.com/cryptowilliam/goutil/container/gbytes"
	"testing"
)

func TestSymmetricEncryptDecrypt(t *testing.T) {
	type TestPack struct {
		alg        Alg
		plainText  []byte
		key        []byte
		wrongKey   []byte
		cipherText []byte
	}

	var pks []TestPack

	pk := TestPack{
		alg:       ALG_AES_128_CBC,
		plainText: []byte("原来是原文啊！原来是原文啊！原来是原文啊！原来是原文啊！原来是原文啊！原来是原文啊！"),
		key:       []byte("mustbe16password"),
		wrongKey:  []byte("mustbe17password"),
	}
	pks = append(pks, pk)

	pk2 := TestPack{
		alg:       ALG_AES_256_CTR,
		plainText: []byte("原来是原文啊！原来是原文啊！原来是原文啊！原来是原文啊！原来是原文啊！原来是原文啊！"),
		key:       []byte("mustbe32passwordmustbe32password"),
		wrongKey:  []byte("mustbe33passwordmustbe33password"),
	}
	pks = append(pks, pk2)

	for _, v := range pks {
		t.Logf("%s, plainText(%s), key(%s)", string(v.alg), string(v.plainText), string(v.key))
		cipherText, err := SymmetricEncrypt(v.alg, v.plainText, v.key)
		if err != nil {
			t.Error(err)
			return
		}
		b58CipherText, _ := gbytes.Encode("base58", cipherText)
		t.Logf("cipherText(%s), base58(%s)", string(cipherText), b58CipherText)
		plainText, err := SymmetricDecrypt(v.alg, cipherText, v.key)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("decrypted plainText(%s)", plainText)
		if !bytes.Equal(plainText, v.plainText) {
			t.Errorf("algorithm %s decrypt error, plainText %s decrypted to %s", v.alg, v.plainText, plainText)
		}
	}
}
