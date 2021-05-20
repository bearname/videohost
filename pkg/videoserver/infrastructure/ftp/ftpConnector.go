package ftp

import (
	"github.com/bearname/videohost/pkg/common/util"
	"github.com/jlaffaye/ftp"
	"io"
)

type FtpConnection struct {
	client *ftp.ServerConn
}

func NewFtpConnection(address string, username string, password string) *FtpConnection {
	client, err := ftp.Dial(address)
	if err != nil {
		return nil
	}

	if err := client.Login(username, password); err != nil {
		return nil
	}
	connection := new(FtpConnection)
	connection.client = client
	return connection
}

func (f *FtpConnection) CopyFile(videoId string, r io.Reader) error {
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
