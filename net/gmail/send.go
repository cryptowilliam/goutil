package gmail

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/sys/gtime"
	gomail "github.com/go-mail/mail"
	"strings"
	"time"
)

// Provider can be null,
// if null, will try to parse built-in provider by From email address.
func Send(e Envelope, c SendContent, password string, p *Provider) error {
	// Validate
	if p == nil {
		err := error(nil)
		p, err = TryParseProvider(e.From.Email)
		if err != nil {
			return err
		}
	}
	if err := p.Validate(); err != nil {
		return err
	}
	if err := e.Validate(); err != nil {
		return err
	}
	if len(password) == 0 {
		return gerrors.Errorf("Send(): Empty password")
	}

	// Convert Envelope to gomail Message
	msg := gomail.NewMessage()
	msg.SetHeader("Subject", e.Subject)
	msg.SetHeader("From", e.From.Email)
	to := []string{}
	for _, v := range e.To {
		to = append(to, msg.FormatAddress(v.Email, v.Showname))
	}
	msg.SetHeader("To", to...)
	cc := []string{}
	for _, v := range e.Cc {
		cc = append(cc, msg.FormatAddress(v.Email, v.Showname))
	}
	msg.SetHeader("Cc", cc...)
	if c.BodyType == BodyTypeHTML {
		msg.SetBody("text/html", c.BodyString)
	} else {
		msg.SetBody("text/plain", c.BodyString)
	}
	for _, v := range c.AttachmentsPath {
		msg.Attach(v)
	}

	// Send the email
	sndAddr, sndPort, ssl, err := p.GetSendServer()
	if err != nil {
		return err
	}
	loginname, err := e.From.LoginName()
	if err != nil {
		return err
	}
	domain, err := e.From.Domain()
	if err != nil {
		return err
	}
	// in ResellerClub free mail account, loginname should be complete email address
	// send address: us2.smtp.mailhostbox.com
	// login name: user@mydomain.com
	if !strings.Contains(sndAddr, domain) {
		loginname = e.From.Email
	}
	d := gomail.NewDialer(sndAddr, sndPort, loginname, password)
	d.SSL = ssl
	if !ssl { // this is required if no ssl/tls/starttls used in the remote mail server
		d.StartTLSPolicy = gomail.NoStartTLS
	}
	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

// send to itself to test it
func TestAccount(addr, pwd string, p *Provider) error {
	evn := Envelope{}
	evn.From.Email = addr
	evn.From.Showname = "email test"
	to := AddrEdit{Email: addr, Showname: ""}
	evn.To = append(evn.To, to)
	evn.Subject = fmt.Sprintf("email test - %s", gtime.Today(time.Local).String())
	content := SendContent{}
	content.BodyString = "email test"
	content.BodyType = BodyTypePlainText
	if err := Send(evn, content, pwd, p); err != nil {
		return err
	}
	return nil
}
