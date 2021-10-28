//go:generate go run gen.go
package sfclassic

import (
	"bytes"
	"compress/flate"
	"io"

	"github.com/richardlehane/siegfried"
	"github.com/richardlehane/siegfried/pkg/core"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Siegfried struct {
	*siegfried.Siegfried
}

func New() *Siegfried {
	rc := flate.NewReader(bytes.NewBuffer(sfcontent))
	sf, err := siegfried.LoadReader(rc)
	check(err)
	check(rc.Close())
	return &Siegfried{sf}
}

func (sf *Siegfried) Identify(rdr io.Reader, name string) ([]core.Identification, error) {
	return sf.Siegfried.Identify(rdr, name, "")
}
