package gcompress

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"strings"
)

type (
	Comp string
)

var (
	allComps []Comp

	CompNone   = enrollComp("none")
	CompSnappy = enrollComp("snappy")
	CompS2      = enrollComp("s2")
	CompZip     = enrollComp("zip")
	CompGzip    = enrollComp("gzip")
	CompPgZip   = enrollComp("pgzip")
	CompZStd    = enrollComp("zstd")
	CompZLib    = enrollComp("zlib")
	CompFlate   = enrollComp("flate")
)

func enrollComp(algo string) Comp {
	allComps = append(allComps, Comp(algo))
	return Comp(algo)
}

func ToComp(algo string) (Comp, error) {
	for _, v := range allComps {
		if string(v) == strings.ToLower(algo) {
			return v, nil
		}
	}
	return CompNone, gerrors.New("unrecognized compress algorithm '%s'", algo)
}