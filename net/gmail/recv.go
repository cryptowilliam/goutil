package gmail

import (
	"encoding/base64"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/i18n/gcharset"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

type IndexType int

const (
	IndexTypeUndefined IndexType = iota + 0
	IndexTypeUid
	IndexTypeSeqNum
)

const (
	INBOX      = "INBOX"
	StartIndex = 1 // First index(Uid or SeqNum) is 1, NOT 0
)

type Filter struct {
	AfterThisTime  *time.Time
	BeforeThisTime *time.Time

	SubjectContains *string
	SenderContains  *string

	IndexType  IndexType
	IndexStart *int
	IndexStop  *int
}

func (f *Filter) Useful() bool {
	if f == nil {
		return false
	}

	if f.AfterThisTime == nil && f.BeforeThisTime == nil &&
		f.SubjectContains == nil && f.SenderContains == nil &&
		f.IndexStart == nil && f.IndexStop == nil {
		return false
	}

	if f.IndexType != IndexTypeUndefined {
		if f != nil && f.IndexStart != nil {
			if *f.IndexStart < StartIndex {
				return false
			}
		}
		if f != nil && f.IndexStop != nil {
			if *f.IndexStop < StartIndex {
				return false
			}
		}
		if f != nil && f.IndexStart != nil && f.IndexStop != nil {
			if *f.IndexStart > *f.IndexStop {
				return false
			}
		}
	}

	return true
}

// Check whether input envelope fits current Filter
func (f *Filter) CheckEnvelope(env *Envelope) bool {
	if f == nil {
		return true
	}
	if env == nil {
		return false
	}

	if f.AfterThisTime != nil && env.recvTime.Before(*f.AfterThisTime) {
		return false
	}
	if f.BeforeThisTime != nil && env.recvTime.After(*f.BeforeThisTime) {
		return false
	}

	if f.SenderContains != nil && strings.Contains(env.From.Email, *f.SenderContains) {
		return false
	}
	if f.SubjectContains != nil && strings.Contains(env.Subject, *f.SubjectContains) {
		return false
	}

	if f.IndexType != IndexTypeUndefined {
		if f.IndexStart != nil && *f.IndexStart > StartIndex {
			if f.IndexType == IndexTypeUid && env.recvUniqueId < uint32(*f.IndexStart) {
				return false
			}
			if f.IndexType == IndexTypeSeqNum && env.recvSeqNum < uint32(*f.IndexStart) {
				return false
			}
		}

		if f.IndexStop != nil && *f.IndexStop > StartIndex {
			if f.IndexType == IndexTypeUid && env.recvUniqueId > uint32(*f.IndexStop) {
				return false
			}
			if f.IndexType == IndexTypeSeqNum && env.recvSeqNum > uint32(*f.IndexStop) {
				return false
			}
		}
	}

	return true
}

func msgToMail(msg *imap.Message, sec *imap.BodySectionName) (*Mail, error) {
	if msg == nil {
		return nil, gerrors.New("Empty source envelope")
	}

	e, err := msgToEnvelop(msg)
	if err != nil {
		return nil, err
	}

	c := RecvContent{}

	r := msg.GetBody(sec)
	if r == nil {
		return nil, gerrors.New("Server didn't returned message body")
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		return nil, err
	}

	// Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, gerrors.Wrap(err, "msgToMail -> NextPart")
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, err := ioutil.ReadAll(p.Body)
			if err != nil {
				return nil, gerrors.Wrap(err, "msgToMail -> ioutil.ReadAll")
			}
			c.BodyString = string(b)
			t, _, err := h.ContentType()
			if err != nil {
				return nil, gerrors.Wrap(err, "msgToMail -> ContentType")
			}
			if t == "text/plain" {
				c.BodyType = BodyTypePlainText
			} else {
				c.BodyType = BodyTypeHTML
			}
		case *mail.AttachmentHeader:
			// This is an attachment
			filename, err := h.Filename()
			if err != nil {
				return nil, gerrors.Wrap(err, "msgToMail -> Filename")
			}
			b, _ := ioutil.ReadAll(p.Body)
			if err != nil {
				return nil, gerrors.Wrap(err, "msgToMail -> AttachmentHeader -> ioutil.ReadAll")
			}
			atch := Attachment{}
			atch.FileName = filename
			atch.Content = b
			c.Attachments = append(c.Attachments, atch)
		}
	}

	return &Mail{Env: *e, Content: c}, nil
}

func decode(s string) string {
	gbk := false

	// Tencent QQ mail returns base64 GBK encoded content
	if gstring.StartWith(s, "=?GBK?B?") && gstring.EndWith(s, "?=") {
		s = gstring.RemoveHead(s, 8)
		s = gstring.RemoveTail(s, 2)
		gbk = true
	}

	decodeBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s
	}
	s = string(decodeBytes)

	if gbk {
		s2, err := gcharset.GbkToUtf8([]byte(s))
		if err == nil {
			s = string(s2)
		}
	}

	return s
}

func msgToEnvelop(msg *imap.Message) (*Envelope, error) {
	if msg == nil {
		return nil, gerrors.New("Empty message")
	}
	if msg == nil || msg.Envelope == nil {
		return nil, gerrors.New("Empty envelope")
	}

	src := msg.Envelope
	res := Envelope{}
	if src.Sender != nil && src.Sender[0] != nil {
		res.From.Email = src.Sender[0].MailboxName
		res.From.Showname = decode(src.Sender[0].PersonalName)
	}

	for _, v := range src.To {
		to := AddrEdit{Email: v.MailboxName, Showname: decode(v.PersonalName)}
		res.To = append(res.To, to)
	}
	for _, v := range src.Cc {
		cc := AddrEdit{Email: v.MailboxName, Showname: decode(v.PersonalName)}
		res.To = append(res.Cc, cc)
	}
	res.Subject = decode(src.Subject)
	res.recvSeqNum = msg.SeqNum
	res.recvUniqueId = msg.Uid
	res.recvTime = msg.InternalDate

	return &res, nil
}

type Receiver struct {
	imap       *client.Client
	currentBox string
}

// IMAP supported, POP3 NOT supported.
//  - p: 可以为空，也可以是指定的provider，可选参数
func NewReceiver(email, password string, p *Provider) (*Receiver, error) {
	// Validate
	if p == nil {
		err := error(nil)
		p, err = TryParseProvider(email)
		if err != nil {
			return nil, err
		}
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}

	// Get login name
	loginname, err := GetLoginName(email)
	if err != nil {
		return nil, err
	}
	domain, err := GetHost(email)
	if err != nil {
		return nil, err
	}
	// in ResellerClub free mail account, loginname should be complete email address
	// send address: us2.smtp.mailhostbox.com
	// login name: user@mydomain.com
	if !strings.Contains(p.IMAPAddress, domain) {
		loginname = email
	}

	// Connect server
	c := new(client.Client)
	if p.IMAPIsSSL {
		c, err = client.DialTLS(p.IMAPAddress, nil)
	} else {
		c, err = client.Dial(p.IMAPAddress)
	}
	if err != nil {
		return nil, err
	}
	if err := c.Login(loginname, password); err != nil {
		return nil, err
	}

	return &Receiver{imap: c}, nil
}

func (s *Receiver) ShowAllBoxes() ([]string, error) {
	doneCh := make(chan error, 1)
	allBoxes := make(chan *imap.MailboxInfo, 10)
	doneCh <- s.imap.List("", "*", allBoxes)
	if err := <-doneCh; err != nil {
		return nil, err
	}

	res := []string{}
	for v := range allBoxes {
		res = append(res, v.Name)
	}
	return res, nil
}

func (s *Receiver) GetBoxInfo(box string) (total, unread int, err error) {
	state, err := s.imap.Select(box, true)
	if err != nil {
		return 0, 0, err
	}
	return int(state.Messages), int(state.Unseen), nil
}

// Notice:
// 如何尽可能少的拉取Envelope: 把起止Uid或者SeqNum作为查询条件放进去
func (s *Receiver) GetEnvelopes(box string, filter *Filter) ([]Envelope, error) {
	// Select box
	ms, err := s.imap.Select(box, true)
	if err != nil {
		return nil, err
	}

	// Set from to tag, in default index is SeqNum
	from := StartIndex
	to := ms.Messages // Mails count in INBOX
	if filter != nil && filter.Useful() {
		if filter.IndexStart != nil {
			from = *filter.IndexStart
		}
		if filter.IndexStop != nil {
			to = uint32(*filter.IndexStop)
		}
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(from), uint32(to))
	// seqset的另一个用法
	// seqset, _ := imap.ParseSeqSet(fmt.Sprintf("%d:", readAfter))

	// Get envelopes after given sequence number
	doneCh := make(chan error, 1)
	res := []Envelope{}
	recvMsgs := make(chan *imap.Message, int(to)-from+1)
	if filter != nil && filter.Useful() && filter.IndexType == IndexTypeUid {
		doneCh <- s.imap.UidFetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid, imap.FetchInternalDate}, recvMsgs)
	} else { // In default, SeqNum accepted. If not set, fetch all envelopes
		doneCh <- s.imap.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid, imap.FetchInternalDate}, recvMsgs)
	}

	// Filter results
	for msg := range recvMsgs {
		evlp, err := msgToEnvelop(msg)
		if err == nil {

			if filter == nil || (filter != nil && filter.CheckEnvelope(evlp)) {
				res = append(res, *evlp)
			}
		}
	}
	if err := <-doneCh; err != nil {
		return nil, err
	}
	return res, nil
}

// Get mails by Uid set
func (s *Receiver) GetMails(box string, uid []int) ([]Mail, error) {
	if uid == nil || len(uid) == 0 {
		return nil, gerrors.New("Invalid input Uid")
	}

	// Select box
	if m2 := s.imap.Mailbox(); m2 == nil || m2.Name != box {
		_, err := s.imap.Select(box, true)
		if err != nil {
			return nil, err
		}
	}

	doneCh := make(chan error, 1)
	seqset := new(imap.SeqSet)
	count := 0
	for _, id := range uid {
		// seqset.AddRange(uint32(uid), uint32(uid))
		seqset.AddNum(uint32(id))
		count++
	}
	recvMsgs := make(chan *imap.Message, count)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope, imap.FetchUid, imap.FetchInternalDate,
		imap.FetchBody, imap.FetchBodyStructure, imap.FetchFlags}
	// 不要使用 s.imap.UidFetch，否则会有少部分邮件读不到正文，使用FetchUid就可以指定Uid读取了
	doneCh <- s.imap.Fetch(seqset, items, recvMsgs)
	if err := <-doneCh; err != nil {
		return nil, err
	}

	if len(recvMsgs) == 0 {
		return nil, gerrors.New("Read null mail")
	}

	res := []Mail{}
	for msg := range recvMsgs {
		m, err := msgToMail(msg, section)
		if err != nil {
			return nil, err
		}
		res = append(res, *m)
	}
	return res, nil
}

func (s *Receiver) Close() error {
	s.imap.Logout()
	return nil
}

func Recv(email, password string, p *Provider, filter *Filter) ([]Mail, error) {
	rcv, err := NewReceiver(email, password, p)
	if err != nil {
		return nil, err
	}
	defer rcv.Close()

	envs, err := rcv.GetEnvelopes(INBOX, filter)
	if err != nil {
		return nil, err
	}
	uids := []int{}
	for _, v := range envs {
		uids = append(uids, int(v.recvUniqueId))
	}

	return rcv.GetMails(INBOX, uids)
}
