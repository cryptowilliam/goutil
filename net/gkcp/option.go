package gkcp

import (
	"crypto/sha1"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

const (
	defaultKey  string = "it's a secrect"
	defaultSalt string = "kcp-go"
)

const (
	modeNormal string = "normal"
	modeFast   string = "fast"
	modeFast2  string = "fast2"
	modeFast3  string = "fast3"
)

// kcp options
type Option struct {
	Key         string `json:"key"`         /* Must be the same between c-s */
	Crypt       string `json:"crypt"`       /* Must be the same between c-s */
	DataShard   int    `json:"datashard"`   /* Must be the same between c-s */
	ParityShard int    `json:"parityshard"` /* Must be the same between c-s */
	MTU         int    `json:"mtu"`
	SndWnd      int    `json:"sndwnd"`
	RcvWnd      int    `json:"rcvwnd"`
	DSCP        int    `json:"dscp"`
	AckNodelay  bool   `json:"acknodelay"`
	SockBuf     int    `json:"sockbuf"` // Shared with smux session
	// mode params
	Mode         string `json:"mode"`
	NoDelay      int    `json:"nodelay"`
	Interval     int    `json:"interval"`
	Resend       int    `json:"resend"`
	NoCongestion int    `json:"nc"`
}

func (opt *Option) getCryptBlock() (kcp.BlockCrypt, error) {
	pass := pbkdf2.Key([]byte(opt.Key), []byte(defaultSalt), 4096, 32, sha1.New)
	var block kcp.BlockCrypt
	err := error(nil)
	switch opt.Mode {
	case "sm4":
		block, err = kcp.NewSM4BlockCrypt(pass[:16])
	case "tea":
		block, err = kcp.NewTEABlockCrypt(pass[:16])
	case "xor":
		block, err = kcp.NewSimpleXORBlockCrypt(pass)
	case "none":
		block, err = kcp.NewNoneBlockCrypt(pass)
	case "aes-128":
		block, err = kcp.NewAESBlockCrypt(pass[:16])
	case "aes-192":
		block, err = kcp.NewAESBlockCrypt(pass[:24])
	case "blowfish":
		block, err = kcp.NewBlowfishBlockCrypt(pass)
	case "twofish":
		block, err = kcp.NewTwofishBlockCrypt(pass)
	case "cast5":
		block, err = kcp.NewCast5BlockCrypt(pass[:16])
	case "3des":
		block, err = kcp.NewTripleDESBlockCrypt(pass[:24])
	case "xtea":
		block, err = kcp.NewXTEABlockCrypt(pass[:16])
	case "salsa20":
		block, err = kcp.NewSalsa20BlockCrypt(pass)
	default:
		block, err = kcp.NewAESBlockCrypt(pass)
	}
	return block, err
}

// Fix no-delay by mode
func (opt *Option) CorrectByMode() {
	switch opt.Mode {
	case modeNormal:
		opt.NoDelay, opt.Interval, opt.Resend, opt.NoCongestion = 0, 40, 2, 1
	case modeFast:
		opt.NoDelay, opt.Interval, opt.Resend, opt.NoCongestion = 0, 30, 2, 1
	case modeFast2:
		opt.NoDelay, opt.Interval, opt.Resend, opt.NoCongestion = 1, 20, 2, 1
	case modeFast3:
		opt.NoDelay, opt.Interval, opt.Resend, opt.NoCongestion = 1, 10, 2, 1
	}
}

// Default is game/Live mode
func DefaultOption(clientSide bool) Option {
	opt := Option{}

	opt.Key = defaultKey
	opt.Crypt = "aes"
	opt.Mode = modeFast3
	opt.MTU = 1350
	if clientSide {
		opt.SndWnd = 128
		opt.RcvWnd = 512
	} else {
		opt.SndWnd = 1024
		opt.RcvWnd = 1024
	}
	opt.DataShard = 10
	opt.ParityShard = 3
	opt.DSCP = 0
	opt.AckNodelay = true
	opt.NoDelay = 0
	opt.Interval = 50
	opt.Resend = 0
	opt.NoCongestion = 0
	opt.SockBuf = 4194304

	opt.CorrectByMode()

	return opt
}

func DefaultSnmp() *kcp.Snmp {
	return kcp.DefaultSnmp
}
