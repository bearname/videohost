package transport

type AddCommentRequest struct {
	VideoId  string `json:"videoId"`
	Message  string `json:"message"`
	ParentId int    `json:"parentId"`
}

type EditCommentRequest struct {
	Message string `json:"message"`
}
