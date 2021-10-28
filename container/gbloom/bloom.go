package gbloom

import (
	"bufio"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/willf/bloom"
	"os"
)

type BloomFilter struct {
	filter        *bloom.BloomFilter
	autoSave      bool
	savePath      string
	writeDiskFlag int
	file          *os.File
}

func New(maxCount uint, maxFp float64, autoSave bool, savePath string) (*BloomFilter, error) {
	if maxCount < 2 || maxFp < 0 || maxFp > 1 {
		return nil, gerrors.New("invalid bloom filter parameters")
	}
	var bf BloomFilter
	bf.autoSave = autoSave
	bf.savePath = savePath
	bf.writeDiskFlag = 0
	var err error
	bf.file, err = os.Open(savePath)
	if err != nil {
		return nil, err
	}

	pi, err := gfs.GetPathInfo(savePath)
	if err != nil {
		return nil, err
	}
	fileSize, err := gfs.GetFileByteSize(savePath)
	if err != nil {
		return nil, err
	}
	if pi.Exist {
		total := int64(0)
		for total < fileSize {
			n, err := bf.filter.ReadFrom(bufio.NewReader(bf.file))
			if err != nil {
				return nil, err
			}
			total += n
		}

	} else {
		bf.filter = bloom.NewWithEstimates(maxCount, maxFp)
	}
	if bf.filter == nil {
		return nil, gerrors.New("bloom filter create fail")
	}

	return &bf, nil
}

func (bf *BloomFilter) Add(data []byte) {
	bf.filter.Add(data)
	bf.writeDiskFlag++
	if bf.autoSave && bf.writeDiskFlag > 100 {
		bf.filter.WriteTo(bufio.NewWriter(bf.file))
		bf.writeDiskFlag = 0
	}
}

func (bf *BloomFilter) AddStr(str string) {
	bf.filter.AddString(str)
	bf.writeDiskFlag++
	if bf.autoSave && bf.writeDiskFlag > 100 {
		bf.filter.WriteTo(bufio.NewWriter(bf.file))
		bf.writeDiskFlag = 0
	}
}

func (bf *BloomFilter) MightContain(data []byte) bool {
	return bf.filter.Test(data)
}

func (bf *BloomFilter) MightContainStr(str string) bool {
	return bf.filter.TestString(str)
}

func (bf *BloomFilter) Reset() error {
	return os.Remove(bf.savePath)
}
