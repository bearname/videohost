package domain

type CommentDto struct {
	UserId   string
	VideoId  string
	Message  string
	ParentId int
}

type CommentService interface {
	Create(commentDto *CommentDto) (int64, error)
	FindRootLevel(videoId string, page *Page) (VideoComments, error)
	FindChildren(rootCommentId int, page *Page) ([]Comment, error)
	Edit(commentId int, message string, userId string) error
	Delete(videoId int, userId string) error
	FindUserComments(id string, page *Page) (Comments, error)
}
