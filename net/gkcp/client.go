package gkcp

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/xtaci/kcp-go"
	"net"
	"time"
)

// kcp client type
type Client struct {
	serverAddr string
	conn       *kcp.UDPSession
	opt        Option
}

type Dialer struct{}

func (d Dialer) Dial(serverAddr string) (net.Conn, error) {
	return DialWithOptions(serverAddr, DefaultOption(true))
}

// Create new client connection and dial
func Dial(serverAddr string) (net.Conn, error) {
	return DialWithOptions(serverAddr, DefaultOption(true))
}

// Create new client connection and dial
func DialWithOptions(serverAddr string, opt Option) (net.Conn, error) {
	var cli Client
	var err error
	cli.opt = opt

	// create a new connection
	block, _ := cli.opt.getCryptBlock()
	cli.conn, err = kcp.DialWithOptions(serverAddr, block, cli.opt.DataShard, cli.opt.ParityShard)
	if err != nil {
		return nil, gerrors.Wrap(err, "createConn()")
	}

	// apply config options
	cli.conn.SetStreamMode(true)
	cli.conn.SetWriteDelay(false)
	cli.conn.SetNoDelay(cli.opt.NoDelay, cli.opt.Interval, cli.opt.Resend, cli.opt.NoCongestion)
	cli.conn.SetWindowSize(cli.opt.SndWnd, cli.opt.RcvWnd)
	cli.conn.SetMtu(cli.opt.MTU)
	cli.conn.SetACKNoDelay(cli.opt.AckNodelay)
	if err := cli.conn.SetDSCP(cli.opt.DSCP); err != nil {
		fmt.Println("SetDSCP:", err)
	}
	if err := cli.conn.SetReadBuffer(cli.opt.SockBuf); err != nil {
		fmt.Println("SetReadBuffer:", err)
	}
	if err := cli.conn.SetWriteBuffer(cli.opt.SockBuf); err != nil {
		fmt.Println("SetWriteBuffer:", err)
	}

	return &cli, nil
}

func (c Client) Read(buf []byte) (n int, err error) {
	return c.conn.Read(buf)
}

func (c Client) Write(data []byte) (n int, err error) {
	return c.conn.Write(data)
}

func (c Client) SetDeadline(t time.Time) error {
	err := c.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}

func (c Client) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c Client) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c Client) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c Client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c Client) Close() error {
	return c.conn.Close()
}
