package gmail

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/grand"
	"github.com/mikaa123/imapmq"
)

// Used email account to implement a simple message queue

type inq struct {
	mq *imapmq.IMAPMQ
	q  *imapmq.Queue
}

type MQ struct {
	inqs []inq
}

type Msg struct {
	Subject string
	Body    []byte
}

func NewConn(as []Account) (*MQ, error) {
	allerr := error(nil)
	res := MQ{}

	// Create new IMAPMQ clients
	for _, a := range as {
		prvd, err := TryParseProvider(a.Email())
		if err != nil {
			allerr = gerrors.New(allerr.Error() + "\r\n" + err.Error())
			continue
		}
		mq, err := imapmq.New(imapmq.Config{
			Login:  a.UserName(),
			Passwd: a.Password(),
			URL:    prvd.IMAPAddress,
		})
		if err != nil {
			allerr = gerrors.New(allerr.Error() + "\r\n" + err.Error())
			continue
		}
		// Create a queue based on INBOX
		q, err := mq.Queue("INBOX")
		if err != nil {
			allerr = gerrors.New(allerr.Error() + "\r\n" + err.Error())
			continue
		}
		res.inqs = append(res.inqs, inq{mq: mq, q: q})
	}

	if len(res.inqs) == 0 {
		return nil, allerr
	}

	return &res, nil
}

func (m *MQ) randomServer() inq {
	return m.inqs[grand.RandomInt(0, len(m.inqs)-1)]
}

func (m *MQ) Pub(subject string, data []byte) {
	m.randomServer().q.Pub(subject, data)
}

/*
func (m *MQ) Sub(subject string, output <-chan *Msg) {
	for _, v := range m.inqs {
		go func() {
			// Subscribe to messages with the subject
			subCh := v.q.Sub(subject)
			for recv := range subCh { // msg is a mail.Message instance.
			buf := bytes.Buffer{}
				buf.ReadFrom(recv.Body)
			item := Msg{Subject:recv.Header.Get("Subject"), Body:buf.Bytes()}
			output <- &item
			}
		}()
	}
}*/

func (m *MQ) Close() error {
	return nil
}
