package util

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

const ContentDir = "content"

const VideoContentType = "video/mp4"
const VideoFileName = "index.mp4"
const ThumbFileName = "screen.jpg"

func CopyFile(fileReader multipart.File, id string) error {
	file, err := CreateFile(VideoFileName, id)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, fileReader)
	if err != nil {
		return err
	}

	return nil
}

func CreateFile(fileName string, id string) (*os.File, error) {
	if err := os.Mkdir(ContentDir, os.ModeDir); err != nil && !os.IsExist(err) {
		return nil, err
	}

	dirPath := filepath.Join(ContentDir, id)
	if err := os.Mkdir(dirPath, os.ModeDir); err != nil && !os.IsExist(err) {
		return nil, err
	}

	filePath := filepath.Join(dirPath, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	return file, err
}
