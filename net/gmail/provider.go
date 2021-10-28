package gmail

import (
	"encoding/json"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/net/gaddr"
	"strings"
)

// Mail service provider
type Provider struct {
	SMTPAddress string
	SMTPIsSSL   bool
	POP3Address string
	POP3IsSSL   bool
	IMAPAddress string
	IMAPIsSSL   bool
	EmailDomain string // Used to parse provider from email address
}

func (p *Provider) Validate() error {
	// Check members are valid
	if len(p.SMTPAddress) == 0 {
		return gerrors.New("Empty SMTP server address")
	}
	if len(p.IMAPAddress) == 0 {
		return gerrors.New("Empty IMAP server address")
	}
	return nil
}

func (p Provider) String() string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

// Get receive server: IMAP / POP3
func (p *Provider) GetReceiveServer() (address string, port int, ssl bool, err error) {
	if len(p.IMAPAddress) > 0 {
		us, err := gaddr.ParseUrl(p.IMAPAddress)
		if err != nil {
			return "", 0, false, gerrors.Errorf("GetRecvServer error for %s", p.String())
		}
		return us.Host.Domain, us.Host.Port, p.IMAPIsSSL, nil
	} else if len(p.POP3Address) > 0 {
		us, err := gaddr.ParseUrl(p.POP3Address)
		if err != nil {
			return "", 0, false, gerrors.Errorf("GetRecvServer error for %s", p.String())
		}
		return us.Host.Domain, us.Host.Port, p.SMTPIsSSL, nil
	} else {
		return "", 0, false, gerrors.Errorf("GetRecvServer error for %s", p.String())
	}
}

// Get send server: SMTP / IMAP
func (p *Provider) GetSendServer() (address string, port int, ssl bool, err error) {
	if len(p.SMTPAddress) > 0 {
		us, err := gaddr.ParseUrl(p.SMTPAddress)
		if err != nil {
			return "", 0, false, gerrors.Errorf("GetSendServer error for %s", p.String())
		}
		return us.Host.Domain, us.Host.Port, p.SMTPIsSSL, nil
	} else if len(p.IMAPAddress) > 0 {
		us, err := gaddr.ParseUrl(p.IMAPAddress)
		if err != nil {
			return "", 0, false, gerrors.Errorf("GetSendServer error for %s", p.String())
		}
		if len(us.Host.Domain) == 0 || us.Host.Port <= 0 {
			return "", 0, false, gerrors.Errorf("GetSendServer error for %s", p.String())
		}
		return us.Host.Domain, us.Host.Port, p.IMAPIsSSL, nil
	} else {
		return "", 0, false, gerrors.Errorf("GetSendServer error for %s", p.String())
	}
}

var (
	InProviders = []Provider{
		{
			SMTPAddress: "smtp.gmail.com:465",
			SMTPIsSSL:   true,
			POP3Address: "pop.gmail.com:995",
			POP3IsSSL:   true,
			IMAPAddress: "imap.gmail.com:993",
			IMAPIsSSL:   true,
			EmailDomain: "gmail.com",
		},
		{
			SMTPAddress: "smtp.qq.com:465",
			SMTPIsSSL:   true,
			POP3Address: "pop.qq.com:995",
			POP3IsSSL:   true,
			IMAPAddress: "imap.qq.com:993",
			IMAPIsSSL:   true,
			EmailDomain: "qq.com",
		},
		{
			SMTPAddress: "smtp.mail.yahoo.com:465",
			SMTPIsSSL:   true,
			POP3Address: "pop.mail.yahoo.com:110",
			POP3IsSSL:   true,
			IMAPAddress: "imap.mail.yahoo.com:993",
			IMAPIsSSL:   true,
			EmailDomain: "yahoo.com",
		},
		{
			SMTPAddress: "smtp.aliyun.com:25",
			SMTPIsSSL:   false,
			POP3Address: "pop3.aliyun.com:110",
			POP3IsSSL:   false,
			IMAPAddress: "imap.aliyun.com",
			IMAPIsSSL:   false,
			EmailDomain: "aliyun.com",
		},
		{
			SMTPAddress: "smtp.163.com:465",
			SMTPIsSSL:   true,
			POP3Address: "pop.163.com:995",
			POP3IsSSL:   true,
			IMAPAddress: "imap.163.com:993",
			IMAPIsSSL:   true,
			EmailDomain: "163.com",
		},
		{
			SMTPAddress: "smtp.126.com:465",
			SMTPIsSSL:   true,
			POP3Address: "pop.126.com:995",
			POP3IsSSL:   true,
			IMAPAddress: "imap.126.com:993",
			IMAPIsSSL:   true,
			EmailDomain: "126.com",
		},
	}
)

func TryParseProvider(email string) (*Provider, error) {
	if err := Validate(email); err != nil {
		return nil, err
	}

	for _, v := range InProviders {
		if gstring.EndWith(strings.ToLower(email), "@"+strings.ToLower(v.EmailDomain)) {
			return &v, nil
		}
	}

	return nil, gerrors.Errorf("Can't find built-in provider for %s", email)
}
