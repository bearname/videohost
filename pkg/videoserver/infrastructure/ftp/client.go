package ftp

import (
	"github.com/bearname/videohost/pkg/common/util"
	"github.com/jlaffaye/ftp"
	"io"
)

type FtpClient struct {
	client *ftp.ServerConn
}

func NewFtpConnection(address string, username string, password string) *FtpClient {
	client, err := ftp.Dial(address)
	if err != nil {
		return nil
	}

	if err := client.Login(username, password); err != nil {
		return nil
	}
	connection := new(FtpClient)
	connection.client = client
	return connection
}

func (f *FtpClient) CopyFile(videoId string, r io.Reader) error {
	err := f.client.MakeDir(videoId)
	if err != nil {
		return err
	}

	filePath := videoId + "\\" + util.VideoFileName
	err = f.client.Stor(filePath, r)
	if err != nil {
		return err
	}

	if err := f.client.Quit(); err != nil {
		return err
	}

	return nil
}

func (f *FtpClient) RemoveDirRecur(path string) error {
	return f.client.RemoveDirRecur(path)
}

func (f *FtpClient) RemoveDir(path string) error {
	return f.client.RemoveDir(path)
}
