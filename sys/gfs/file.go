package gfs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gvolume"
	"github.com/cryptowilliam/goutil/sys/gio"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// copy file and notify copied size
func CopyFileEx(src, dst string, sizeCallback gio.CopiedSizeCallback) (size int64, err error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	written, err := gio.CopyEx(out, nil, in, nil, time.Duration(0), sizeCallback)
	if err != nil {
		return written, err
	}

	err = out.Sync()
	if err != nil {
		return written, err
	}

	si, err := os.Stat(src)
	if err != nil {
		return written, err
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return written, err
	}

	return written, nil
}

func FilenameToReader(filename string) (io.Reader, *os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	return io.Reader(file), file, nil
}

func FilenameToWriter(filename string) (io.Writer, *os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	return io.Writer(file), file, nil
}

func FileToBytes(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func FileToString(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	return string(b), err
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func FileToJson(filename string, ptrJsonStruct interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, ptrJsonStruct)
}

func JsonToFile(jsonStruct interface{}, indent bool, filename string) error {
	if jsonStruct == nil {
		return gerrors.New("Null input jsonStruct")
	}
	var b []byte
	err := error(nil)
	if indent {
		b, err = json.MarshalIndent(jsonStruct, "", "\t")
	} else {
		b, err = json.Marshal(jsonStruct)
	}
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(b)
	return nil
}

func BytesToFile(data []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	total := len(data)
	done := 0
	for {
		if done >= total {
			break
		}
		n, err := f.Write(data)
		if err != nil {
			return err
		}
		done += n
	}
	return nil
}

func StringToFile(data string, filename string) error {
	return BytesToFile([]byte(data), filename)
}

func AppendStringToFile(data string, filename string) error {
	return AppendBytesToFile([]byte(data), filename)
}

func CreateAppendFile(path string) (*os.File, error) {
	// 0660 is okay, don't use os.ModeAppend instead, otherwise there will be "permission not allowed" error when non-admin open the file
	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
}

func AppendBytesToFile(data []byte, filename string) error {
	// 0660 is okay, don't use os.ModeAppend instead, otherwise there will be "permission not allowed" error when non-admin open the file
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(data); err != nil {
		return err
	}
	return nil
}

func GetFileByteSize(filename string) (int64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func GetFileSize(filename string) (gvolume.Volume, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return gvolume.FromByteSize(float64(fi.Size()))
}

/*
// FIXME: can't cross compile for windows under macOS
func MoveToTrash(filename string) error {
	_, err := trash.MoveToTrash(filename)
	return err
}*/
