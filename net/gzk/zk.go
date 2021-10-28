package gzk

import (
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"time"
)

type (
	ZK struct {
		conn *zk.Conn
	}
)

func Dial(servers []string, timeout time.Duration) (*ZK, error) {
	conn, _, err := zk.Connect(servers, timeout)
	if err != nil {
		fmt.Println("connect error", servers)
		return nil, err
	}
	return &ZK{conn: conn}, nil
}

func (zk *ZK) Ls(path string) ([]string, error) {
	res, _, err := zk.conn.Children(path)
	return res, err
}

func (zk *ZK) GetAddr(path string) (string, error) {
	buf, _, err := zk.conn.Get(path)
	if err != nil {
		return "", err
	}
	addr := struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}{}
	fmt.Println(string(buf))
	if err := json.Unmarshal(buf, &addr); err != nil {
		return "", err
	}
	if addr.Address == "" || addr.Port == 0 {
		return "", gerrors.New("address parse error")
	}
	return addr.Address + ":" + strconv.Itoa(addr.Port), nil
}

func (zk *ZK) Close() {
	zk.conn.Close()
}
