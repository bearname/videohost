package domain

import "time"

type BaseComment struct {
	Id      int
	UserId  string
	VideoId string
	Created time.Time
	Message string
}

type RootComment struct {
	BaseComment
	CountSubComments int
}

type Comment struct {
	BaseComment
	ParentId int
}

func NewComment(userId string, videoId string, parentId int, message string) *Comment {
	c := new(Comment)
	c.UserId = userId
	c.VideoId = videoId
	c.Message = message
	c.ParentId = parentId
	return c
}

type Comments struct {
	CountAllComments int
	RootComments     []RootComment
}

type VideoComments struct {
	Comments
	VideoId string
}

func NewVideoComments(videoId string, countAllComments int, rootComments []RootComment) *VideoComments {
	v := new(VideoComments)
	v.VideoId = videoId
	v.RootComments = rootComments
	v.CountAllComments = countAllComments
	return v
}

type Page struct {
	Size   int
	Number int
}

type CommentRepository interface {
	Create(*Comment) (int64, error)
	FindById(commentId int) (Comment, error)
	FindRootLevel(videoId string, page *Page) (VideoComments, error)
	FindChildren(rootCommentId int, page *Page) ([]Comment, error)
	FindUserComments(userId string, page *Page) (Comments, error)
	Edit(commentId int, message string) error
	Delete(commentId int) error
}
