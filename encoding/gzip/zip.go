package gzip

import (
	"archive/zip"
	"bytes"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ZipFile(srcDir, destFilename string) error {
	return nil
}

func UnZipFile(srcFilename, dstDir string) error {
	r, err := zip.OpenReader(srcFilename)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dstDir, os.ModePerm)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dstDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func ZipBuf(srcBuf []byte) ([]byte, error) {
	return nil, nil
}

func UnZip(r io.ReadCloser, outputDir string) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return err
	}

	// Read all the files from zip archive
	for _, v := range zipReader.File {
		if v.FileInfo().IsDir() {
			os.MkdirAll(filepath.Join(outputDir, v.Name), os.ModePerm)
		} else {
			// unzippedFileBytes is unzipped file bytes
			unzippedFileBytes, err := readZipFile(v)
			if err != nil {
				continue
			}
			gfs.BytesToFile(unzippedFileBytes, filepath.Join(outputDir, v.Name))
		}
	}

	return nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
