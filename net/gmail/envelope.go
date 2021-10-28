package gmail

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"strings"
	"time"
)

// 专用于Envelope中用户填写方便而准备的结构体
type AddrEdit struct {
	Email    string
	Showname string
}

func (ae *AddrEdit) Validate() error {
	return Validate(ae.Email)
}

func (ae *AddrEdit) LoginName() (string, error) {
	a, err := NewAddress(ae.Email, ae.Showname)
	if err != nil {
		return "", err
	}
	return a.UserName(), nil
}

func (ae *AddrEdit) Domain() (string, error) {
	ss := strings.Split(ae.Email, "@")
	if len(ss) != 2 {
		return "", gerrors.Errorf("invalid email address(%s)", ae.Email)
	}
	return ss[1], nil
}

type BodyType int

const (
	BodyTypePlainText BodyType = iota + 0
	BodyTypeHTML
)

type Envelope struct {
	Datetime     time.Time
	Subject      string
	From         AddrEdit
	To           []AddrEdit
	Cc           []AddrEdit
	recvUniqueId uint32 // 收邮件时才会有的唯一Id
	recvSeqNum   uint32 // 收邮件时才会有的邮件序号，每封邮件的序号不是固定的，即使部分邮件被用户删除了，序号也是连续的
	recvTime     time.Time
	/*Body                  string
	BodyType BodyType
	Attachment            []string // File path in disk*/
}

type SendContent struct {
	BodyType        BodyType
	BodyString      string
	AttachmentsPath []string
}

type Attachment struct {
	FileName string
	Content  []byte
}

type RecvContent struct {
	BodyType    BodyType
	BodyString  string
	Attachments []Attachment
}

type Mail struct {
	Env     Envelope
	Content RecvContent
}

func (e *Envelope) Validate() error {
	if err := e.From.Validate(); err != nil {
		return err
	}
	for _, v := range e.To {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	for _, v := range e.Cc {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Envelope) Uid() int {
	return int(e.recvUniqueId)
}
