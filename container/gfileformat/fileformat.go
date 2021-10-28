package gfileformat

// github.com/richardlehane/siegfried is a great and professional go library to detect file format
// much more better than any known library

import (
	"fmt"
	"github.com/cryptowilliam/goutil/container/gvolume"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/gabriel-vasile/mimetype"
	"github.com/h2non/filetype"
	"github.com/ross-spencer/sfclassic"
)

// 不太准，比如ttf、Log4j2Plugins.dat都被认为是unix平台的可执行程序
func IsExecutableUnstable(filename string) (bool, error) {

	mime, err := mimetype.DetectFile(filename)
	if err != nil {
		return false, err
	}
	fmt.Println(fmt.Sprintf("%s, mime:%s.", filename, mime))
	if mime.String() == "application/x-executable" || // linux executable file
		mime.String() == "application/octet-stream" || // unix(including darwin, plan9) executable file
		mime.String() == "application/vnd.microsoft.portable-executable" { // windows executable file
		return true, nil
	}

	// linux executable file
	if mime.String() == "application/x-executable" {
		return true, nil
	}

	// windows executable file
	if mime.String() == "application/vnd.microsoft.portable-executable" {
		return true, nil
	}

	// unix(including darwin, plan9) executable file
	if mime.String() == "application/octet-stream" {
		return true, nil
	}

	return false, nil
}

// 不太准，哪些不准忘记了，反正肯定不太准
func IsExecutableUnstable2(filename string) (bool, error) {
	/*sz, err := xfs.GetFileSize(filename)
	if err != nil {
		return false, err
	}
	if sz > errIfBiggerThan {
		return false, gerrors.Errorf("%s it bigger than %s", filename, errIfBiggerThan.String())
	}*/
	b, err := gfs.FileToBytes(filename)
	if err != nil {
		return false, err
	}
	t, err := filetype.Get(b)
	if err != nil {
		return false, err
	}
	fmt.Println(filename, t.MIME, t.Extension)
	if t.Extension == "exe" || t.Extension == "elf" {
		return true, nil
	}
	return false, nil
}

// WARN: 无法识别plan9平台下的可执行程序，但已经是最目前靠谱的识别了
func IsExecutable(filename string) (bool, error) {
	ff, err := DetectFileFormat(filename)
	if err != nil {
		if vol, err2 := gfs.GetFileSize(filename); err2 == nil && vol == gvolume.Volume(0) {
			return false, nil
		}
		return false, err
	}
	if ff.Name == "Mach-O" ||
		ff.Name == "Executable and Linkable Format" ||
		ff.Name == "Windows Portable Executable" {
		return true, nil
	}
	return false, nil
}

type FileFormat struct {
	NameSpace string
	Name      string
	Version   string
	MIME      string
}

// WARN: 无法识别plan9平台下的可执行程序，但已经是最目前靠谱的识别了
// DetectFileFormat返回了多个命名空间下的识别结果，而sf程序只选取了pronom命名空间
// if file is empty, returns error
func DetectFileFormat(filename string) (*FileFormat, error) {
	rd, fd, err := gfs.FilenameToReader(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	ids, err := sfclassic.New().Identify(rd, "classic.sig")
	if err != nil {
		return nil, err
	}

	var r []FileFormat
	for _, v := range ids {
		if len(v.Values()) == 7 {
			r = append(r, FileFormat{NameSpace: v.Values()[0], Name: v.Values()[2], Version: v.Values()[3], MIME: v.Values()[4]})
		}
	}
	if len(r) == 0 {
		return nil, gerrors.Errorf("detect failed")
	}
	return &r[0], nil
}
