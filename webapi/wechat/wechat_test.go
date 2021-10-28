package wechat

import (
	"fmt"
	"testing"
)

func TestWechat(t *testing.T) {
	recvch := make(chan RecvMsg, 100)
	s, err := NewBot(LoginModeWeb, recvch, nil)
	if err != nil {
		t.Error(err)
		return
	}

	myself := s.GetMyself()
	fmt.Println("myself", myself)

	users := s.GetAllMyContacts()
	for _, v := range users {
		fmt.Println("myfriend", v)
	}
}
