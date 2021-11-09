package gcompress

type (
	CompAlgo string
)

var (
	CompAlgoSnappy CompAlgo = "snappy"
	CompAlgoS2     CompAlgo = "s2"
	CompAlgoZip    CompAlgo = "zip"
	CompAlgoGzip   CompAlgo = "gzip"
	CompAlgoPgZip  CompAlgo = "pgzip"
	CompAlgoZStd   CompAlgo = "zstd"
	CompAlgoZLib   CompAlgo = "zlib"
	CompAlgoFlate  CompAlgo = "flate"
)
