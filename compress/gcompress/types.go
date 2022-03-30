package gcompress

type (
	Comp string
)

var (
	CompNone   Comp = "none"
	CompSnappy Comp = "snappy"
	CompS2     Comp = "s2"
	CompZip    Comp = "zip"
	CompGzip   Comp = "gzip"
	CompPgZip  Comp = "pgzip"
	CompZStd   Comp = "zstd"
	CompZLib   Comp = "zlib"
	CompFlate  Comp = "flate"
)
