package gfes

import (
	"bytes"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

func defaultTestCases() *gtest.CaseList {
	cl := gtest.NewCaseList()
	cl.New().Input([]byte("hello world")).Input("user secret").Input("salt secret")
	cl.New().Input([]byte("oishi")).Input("t6547654").Input("fesf6576nhh45dm3vFRH")
	cl.New().Input([]byte("美丽")).Input("93w6rtRRRt").Input("glo24r54EF32SDHR4tgtk")
	cl.New().Input([]byte("great")).Input("%#(&#$^&").Input("gggERFJMVGGvffgrh1ad5")
	cl.New().Input([]byte("Wonderful World")).Input("``").Input("p;jhgffdbng$^^$W#")
	cl.New().Input([]byte("Captain Nemo")).Input("vffCC").Input("                 ")
	return cl
}

func TestMartolodDecrypt(t *testing.T) {
	for _, v := range defaultTestCases().Get() {
		plain := v.Inputs[0].([]byte)
		userSecret := v.Inputs[1].(string)
		saltSecret := v.Inputs[2].(string)
		cipher, err := MartolodEncrypt(plain, userSecret, saltSecret)
		gtest.Assert(t, err)
		decPlain, err := MartolodDecrypt(cipher, userSecret, saltSecret)
		gtest.Assert(t, err)
		if !bytes.Equal(plain, decPlain) {
			gtest.Assert(t, gerrors.New("plain %s != decoded plain %s", plain, decPlain))
		}
	}
}

func TestTriMartolodDecrypt(t *testing.T) {
	for _, v := range defaultTestCases().Get() {
		plain := string(v.Inputs[0].([]byte))
		userSecret := v.Inputs[1].(string)
		saltSecret := v.Inputs[2].(string)
		cipher, err := TriMartolodEncrypt(plain, userSecret, saltSecret)
		gtest.Assert(t, err)
		decPlain, err := TriMartolodDecrypt(cipher, userSecret, saltSecret)
		gtest.Assert(t, err)
		if plain != decPlain {
			gtest.Assert(t, gerrors.New("plain %s != decoded plain %s", plain, decPlain))
		}
	}
}

func TestSonnefesEncrypt(t *testing.T) {
	in := "hello 你好 こんにちは"
	key := "this is a very complex key!"
	cipher, err := SonnefesEncrypt(in, key)
	if err != nil {
		t.Error(err)
		return
	}
	dec, err := SonnefesDecrypt(cipher, key)
	if err != nil {
		t.Error(err)
		return
	}
	if dec != in {
		t.Errorf("SonnefesEncrypt/SonnefesDecrypt error")
		return
	}
}