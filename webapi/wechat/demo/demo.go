package main

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/webapi/wechat"
	"time"
)

func main() {
	recvch := make(chan wechat.RecvMsg, 100)
	s, err := wechat.NewBot(wechat.LoginModeTerminal, recvch, glog.DefaultLogger)
	if err != nil {
		glog.Erro(err)
		return
	}

	if err := s.WaitLogin(); err != nil {
		glog.Erro(err)
		return
	}

	users := s.GetAllMyContacts()
	for _, v := range users {
		fmt.Println("myfriend", v)
	}

	myself := s.GetMyself()
	fmt.Println("myself", myself)

	s.SendTextQueue(wechat.RecvNameTypeNickName, "tom", "你🐒啊")

	s.SendTextQueue(wechat.RecvNameTypeNickName, "机器人", "大嘎🐒啊")

	for {
		time.Sleep(time.Second)
	}

}
