package gssh

// Another ssh wrapper which is more powerful but ssh only, no sftp support:
// https://github.com/yahoo/vssh

import (
	"github.com/cryptowilliam/goutil/net/gaddr"
	"github.com/melbahja/goph"
)

type Client struct {
	in *goph.Client
}

// Dial SSH server.
// privateKeyFile: private file path like .pem or .id_rsa
func Dial(address, username, password, privateKeyFile, passphrase string) (*Client, error) {
	// Build auth.
	var auth goph.Auth
	if password != "" {
		auth = goph.Password(password)
	} else {
		err := error(nil)
		auth, err = goph.Key(privateKeyFile, passphrase)
		if err != nil {
			return nil, err
		}
	}

	// Parse address.
	us, err := gaddr.ParseUrl(address)
	if err != nil {
		return nil, err
	}

	// Start new ssh connection with private key.
	res := &Client{}
	res.in, err = goph.New(username, us.Host.String(), auth)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Execute your command and get your output as string.
func (c *Client) RunCommand(cmd string) (string, error) {
	out, err := c.in.Run(cmd)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Upload file with sftp.
func (c *Client) UploadFile(localFile, remoteFile string) error {
	return c.in.Upload(localFile, remoteFile)
}

// Download file with sftp.
func (c *Client) DownloadFile(remoteFile, localFile string) error {
	return c.in.Download(remoteFile, localFile)
}
