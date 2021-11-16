package gfes

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/crypto/gbase"
	"github.com/cryptowilliam/goutil/crypto/gencrypt"
	"github.com/cryptowilliam/goutil/crypto/ghash"
	"github.com/jbenet/go-base58"
)

// Note: DON'T DELETE, DON'T MODIFY CODE
// My Custom Encryption Standard
//
// SiliconHash:
// secret required string hash, when perform multiple times, it returns same output for unchanging input.
// Used to protect the plain text in the process of comparing passwords.
// It is a hash function, it is not reversible.
//
// Sonnefes:
// string to string encrypt/decrypt, used to protect password in PlainText like json/toml config file,
// it provides confidentiality but NO integrity.
//
// Martolod:
// bytes to bytes encrypt/decrypt, used to protect data on the cloud,
// it provides confidentiality and integrity out of the box.
// Note:
// When perform multiple times, it returns different results even user don't change input,
// because random number used in encrypt process, please DON'T use it when you need same output for unchanging input.
//
// TriMartolod:
// string to string encrypt/decrypt based on 'Martolod',
// used to protect password/api-secret in PlainText like json/toml config file,
// it provides confidentiality and integrity out of the box.
// Note:
// When perform multiple times, it returns different results even user don't change input,
// because random number used in encrypt process, please DON'T use it when you need same output for unchanging input.

const (
	head = "1024"
)

func SiliconHash(key string, secret string) string {
	return gbase.Base64Encode(ghash.GetHMAC(ghash.HashTypeSHA512, []byte(key), []byte(secret)))
}

func SonnefesEncrypt(plainText, key string) (string, error) {
	keySHA256 := ghash.GetSHA256([]byte(key))
	cipherAES256, err := gencrypt.SymmetricEncrypt(gencrypt.ALG_AES_256_CBC, []byte(plainText), keySHA256)
	if err != nil {
		return "", err
	}
	return base58.Encode(cipherAES256), nil
}

func SonnefesDecrypt(cipher string, key string) (string, error) {
	keySHA256 := ghash.GetSHA256([]byte(key))
	bin := base58.Decode(cipher)
	plainText, err := gencrypt.SymmetricDecrypt(gencrypt.ALG_AES_256_CBC, bin, []byte(keySHA256))
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func MartolodEncrypt(plain []byte, userSecret, saltSecret string) ([]byte, error) {
	if len(plain) == 0 {
		return nil, gerrors.New("empty plain")
	}
	if userSecret == "" {
		return nil, gerrors.New("empty user secret")
	}
	if saltSecret == "" {
		return nil, gerrors.New("empty salt secret")
	}

	secretHMAC := []byte(userSecret)
	for i := 0; i < 6; i++ {
		secretHMAC = ghash.GetHMAC(ghash.HashTypeSHA256, secretHMAC, []byte(head+saltSecret))
	}

	cipher, err := gencrypt.NewAesGcm256().Encrypt(plain, secretHMAC, true)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}

func MartolodDecrypt(cipher []byte, userSecret, saltSecret string) ([]byte, error) {
	if len(cipher) == 0 {
		return nil, gerrors.New("empty cipher")
	}
	if userSecret == "" {
		return nil, gerrors.New("empty user secret")
	}
	if saltSecret == "" {
		return nil, gerrors.New("empty salt secret")
	}

	secretHMAC := []byte(userSecret)
	for i := 0; i < 6; i++ {
		secretHMAC = ghash.GetHMAC(ghash.HashTypeSHA256, secretHMAC, []byte(head+saltSecret))
	}

	plain, err := gencrypt.NewAesGcm256().Decrypt(cipher, secretHMAC, true)
	if err != nil {
		return nil, err
	}
	return plain, nil
}

func TriMartolodEncrypt(plain string, userSecret, saltSecret string) (string, error) {
	buf, err := MartolodEncrypt([]byte(plain), userSecret, saltSecret)
	if err != nil {
		return "", err
	}
	return gbase.Base64Encode(buf), nil
}

func TriMartolodDecrypt(cipher string, userSecret, saltSecret string) (string, error) {
	buf, err := gbase.Base64Decode(cipher)
	if err != nil {
		return "", err
	}
	res, err := MartolodDecrypt(buf, userSecret, saltSecret)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
