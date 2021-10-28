package gcmd

import "os"

func ChmodAddX(filename string) error {
	return os.Chmod(filename, os.ModePerm)
}
