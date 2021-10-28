package gsecureconn

import (
	"net"
)

type Listener struct {
	lis    net.Listener
	cipher *Cipher
}

func WrapListener(lis net.Listener, cipher *Cipher) net.Listener {
	return &Listener{
		lis:    lis,
		cipher: cipher,
	}
}

// see net.ListenTCP
func (l *Listener) Accept() (net.Conn, error) {
	localConn, err := l.lis.Accept()
	if err != nil {
		return nil, err
	}
	// localConn被关闭时直接清除所有数据 不管没有发送的数据
	//localConn.SetLinger(0)
	return &SecureConn{
		conn:   localConn,
		cipher: l.cipher,
	}, nil
}

func (l *Listener) Addr() net.Addr {
	return l.lis.Addr()
}

func (l *Listener) Close() error {
	return l.Close()
}
