package dto

import (
	"mime/multipart"
)

type UploadVideoDto struct {
	Title         string
	Description   string
	MultipartFile multipart.File
	FileHeader    *multipart.FileHeader
	Chapters      []ChapterDto
}
