package netdrive

import "github.com/cryptowilliam/goutil/basic/gerrors"

type (
	Provider string

	NetDrive interface {
		DownloadFile(pathToFile string) ([]byte, error)
		UploadFile(pathToFile string, buf []byte) error
		Close() error
	}
)

const (
	GoogleDrive Provider = "google-drive"
)

func New(p Provider, apiKey string) (NetDrive, error) {
	switch p {
	case GoogleDrive:
		return newGoogleDrive(apiKey)
	}
	return nil, gerrors.New("unknown network dirve provider %s", p)
}
