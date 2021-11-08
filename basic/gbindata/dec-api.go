package gbindata

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func Dec(fileHexString *string, outputBinaryFilename string) error {
	if fileHexString == nil {
		return gerrors.Errorf("fileHexString is nil")
	}
	if len(*fileHexString)%2 != 0 {
		return gerrors.Errorf("fileHexString length error")
	}
	buf, err := hexutil.Decode(*fileHexString)
	if err != nil {
		return err
	}
	return gfs.BytesToFile(buf, outputBinaryFilename)
}
