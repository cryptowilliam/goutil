package gmail

import (
	"github.com/badoux/checkmail"
	"github.com/cryptowilliam/goutil/container/gstring"
	"strings"
)

type Address struct {
	email string
	//loginname string
	showname string
}

func Validate(email string) error {
	return checkmail.ValidateFormat(email)
}

func GetLoginName(email string) (string, error) {
	addr, err := NewAddress(email, "")
	if err != nil {
		return "", err
	}
	return addr.UserName(), nil
}

func GetHost(email string) (string, error) {
	addr, err := NewAddress(email, "")
	if err != nil {
		return "", err
	}
	return addr.Host(), nil
}

func NewAddress(email, showname string) (*Address, error) {
	if err := Validate(email); err != nil {
		return nil, err
	}
	/*atidx := xstring.IndexAfter(email, "@", 0)
	loginname := email[0:atidx]*/
	return &Address{email: email, showname: showname}, nil
}

func (a *Address) Email() string {
	return a.email
}

func (a *Address) EmailReplaceLoginNameTail(replace string) string {
	return gstring.TrySubstrLenAscii(a.UserName(), 0, len(a.UserName())-len(replace)) + replace + "@" + a.Host()
}

// in quant.lol, login name is username@quant.lol, not username only
func (a *Address) UserName() string {
	ss := strings.Split(a.email, "@")
	return ss[0]
	//atidx := xstring.IndexAfter(a.email, "@", 0)
	//return a.email[0:atidx]
}

func (a *Address) ShowName() string {
	return a.showname
}

func (a *Address) Host() string {
	ss := strings.Split(a.email, "@")
	return ss[1]
	//return xstring.LastSubstrByLenAscii(a.email, len(a.email) - (len(a.loginname) + len("@")))
}
