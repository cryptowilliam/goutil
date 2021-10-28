package deepcopy

import (
	"github.com/ulule/deepcopier"
)

func Copy(src interface{}, dst interface{}) error {
	return deepcopier.Copy(src).To(dst)
}
