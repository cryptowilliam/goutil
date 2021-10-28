package gmail

import "testing"

func TestRecv(t *testing.T) {
	es, err := Recv("", "", nil, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(es)
}
