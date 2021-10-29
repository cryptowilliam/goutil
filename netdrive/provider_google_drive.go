package netdrive

import (
	"context"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"strings"
)

type (
	googleDrive struct {
		srv                *drive.Service
		sharedRootFolder   string
		sharedRootFolderID string
	}
)

func newGoogleDrive(apiKey string) (*googleDrive, error) {
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, gerrors.New("Unable to retrieve drive Client %v", err)
	}
	return &googleDrive{srv: srv}, nil
}

func (gd *googleDrive) DownloadFile(pathToFile string) ([]byte, error) {
	return nil, nil
}

func (gd *googleDrive) UploadFile(pathToFile string, buf []byte) error {
	return nil
}

func (gd *googleDrive) CreateFolder(name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := gd.srv.Files.Create(d).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (gd *googleDrive) buildQuerySearchFile(filePaths []string) string {
	escapedFilePath := make([]string, 0)
	for _, filePath := range filePaths {
		sliceFilePath := strings.Split(filePath, "'")
		escapedFilePath = append(escapedFilePath, strings.Join(sliceFilePath, `\'`))
	}
	subQueries := make([]string, 0)
	for _, file := range escapedFilePath[0:(len(escapedFilePath) - 1)] {
		subQuery := fmt.Sprintf("(name = '%s' and mimeType = 'application/vnd.google-apps.folder')", file)
		subQueries = append(subQueries, subQuery)
	}
	lastQuery := fmt.Sprintf("(name = '%s' and mimeType != 'application/vnd.google-apps.folder')", escapedFilePath[len(escapedFilePath)-1])
	subQueries = append(subQueries, lastQuery)
	return strings.Join(subQueries, " or ")
}

// Return fileID, isExisted, error
func (gd *googleDrive) getFileIDByPath(filePath string) (string, bool, error) {
	listFileInPath := strings.Split(filePath, "/")
	numPathLevel := len(listFileInPath)
	query := gd.buildQuerySearchFile(listFileInPath)
	fileList, err := gd.srv.Files.List().Fields("files(id, name, parents)").Q(query).Do()
	if err != nil {
		return "", false, err
	}
	if len(fileList.Files) < numPathLevel {
		return "", false, nil
	}

	lastFileID := ""
	isExistedPath := true
	preNodeID := gd.sharedRootFolderID
	for _, fileInPath := range listFileInPath[1:] {
		existedNode := false
		for _, file := range fileList.Files {
			if file.Name == fileInPath && file.Parents[0] == preNodeID {
				preNodeID = file.Id
				lastFileID = file.Id
				existedNode = true
				break
			}
		}
		if existedNode == false {
			isExistedPath = false
			break
		}
	}
	if isExistedPath {
		return lastFileID, true, nil
	}
	return "", false, nil
}

func (gd *googleDrive) getReaderByFilePath(filePath string) (string, io.Reader, error) {
	id, existed, err := gd.getFileIDByPath(filePath)
	if err != nil {
		return "", nil, err
	}
	if !existed {
		return "", nil, gerrors.ErrNotExist
	}
	stream, err := gd.getReaderByID(id)
	return id, stream, err
}

func (gd *googleDrive) isExistedByID(id string) (bool, error) {
	files, err := gd.srv.Files.List().Do()
	if err != nil {
		return false, err
	}
	if len(files.Files) == 0 {
		return false, err
	}
	for _, file := range files.Files {

		if file.Id == id {
			return true, nil
		}
	}
	return false, nil
}

func (gd *googleDrive) getReaderByID(fileOrFolderID string) (io.Reader, error) {
	resp, err := gd.srv.Files.Get(fileOrFolderID).Download()
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (n *googleDrive) Close() error {
	return nil
}
