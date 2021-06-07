package ftp

import (
	"github.com/bearname/videohost/pkg/common/util"
	"github.com/jlaffaye/ftp"
	"io"
)

type Client struct {
	client *ftp.ServerConn
}

func NewFtpConnection(address string, username string, password string) *Client {
	client, err := ftp.Dial(address)
	if err != nil {
		return nil
	}

	if err = client.Login(username, password); err != nil {
		return nil
	}
	connection := new(Client)
	connection.client = client
	return connection
}

func (f *Client) CopyFile(videoId string, r io.Reader) error {
	err := f.client.MakeDir(videoId)
	if err != nil {
		return err
	}

	filePath := videoId + "\\" + util.VideoFileName
	err = f.client.Stor(filePath, r)
	if err != nil {
		return err
	}

	return f.client.Quit()
}

func (f *Client) RemoveDirRecur(path string) error {
	return f.client.RemoveDirRecur(path)
}

func (f *Client) RemoveDir(path string) error {
	return f.client.RemoveDir(path)
}
