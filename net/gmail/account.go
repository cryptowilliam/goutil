package gmail

import "github.com/cryptowilliam/goutil/basic/gerrors"

type Account struct {
	addr     Address
	password string
}

func NewAccount(email, showname, password string) (*Account, error) {
	addr, err := NewAddress(email, showname)
	if err != nil {
		return nil, err
	}
	if len(password) == 0 {
		return nil, gerrors.New("Empty password")
	}
	return &Account{addr: *addr, password: password}, nil
}

func (a *Account) Email() string {
	return a.addr.email
}

func (a *Account) UserName() string {
	return a.addr.UserName()
}

func (a *Account) ShowName() string {
	return a.addr.showname
}

func (a *Account) Password() string {
	return a.password
}
