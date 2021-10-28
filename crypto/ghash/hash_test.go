package ghash

import (
	"fmt"
	"testing"
)

func TestGetHMAC(t *testing.T) {
	plain := []byte("hello hmac")
	secret := []byte("it is a secret")
	cipher := GetHMAC(HashTypeSHA512, plain, secret)
	fmt.Println(string(cipher))
}
