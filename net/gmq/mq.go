package gmq

type (
	SubCallback   func(subj string, b []byte)
	ReplyCallback func(service string) ([]byte, error)

	MQ interface {
		Pub(subj string, b []byte) error
		Sub(subj string, callback SubCallback) error
		UnSub(subj string, callback SubCallback) error

		Push(queue string, b []byte) error
		Pop(queue string) ([]byte, error)

		Request(service string, b []byte) error
		Reply(service string, callback ReplyCallback) error
		UnReply(service string, callback ReplyCallback) error
	}
)
