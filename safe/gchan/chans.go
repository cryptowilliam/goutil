package gchan

import "github.com/cryptowilliam/goutil/basic/gerrors"

// reference: https://go101.org/article/channel-closing.html

var errClosedChan = gerrors.New("closed chan")

// Safe close chan with wrap function
func SafeCloseChanStruct(ch chan struct{}) (err error) {
	defer func() {
		if recover() != nil {
			err = errClosedChan
		}
	}()

	// assume ch != nil here.
	if ch != nil {
		close(ch) // panic if ch is closed
	}
	return nil
}

func SafeSendChanStruct(ch chan struct{}, value struct{}) (err error) {
	defer func() {
		if recover() != nil {
			err = errClosedChan
		}
	}()

	if ch == nil {
		return errClosedChan
	}
	ch <- value // panic if ch is closed
	return nil  // <=> closed = false; return
}
