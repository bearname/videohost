package transport

import (
	"github.com/bearname/videohost/internal/common/db"
	"net/http"
	"strconv"
)

const defaultPageSize = 15
const NotSet = -1

type CommentsFilter struct {
	Page     db.Page
	RootId   int
	AuthorId string
	VideoId  string
}

func DecodeFindCommentsRequest(r *http.Request) (CommentsFilter, error) {
	query := r.URL.Query()
	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = defaultPageSize
	}
	pageNum, err := strconv.Atoi(query.Get("page_num"))
	if err != nil || pageNum < 0 {
		pageNum = 0
	}
	pageSpec := &db.Page{
		Size:   pageSize,
		Number: pageNum,
	}
	var rootId int

	rootId, err = strconv.Atoi(query.Get("rootId"))
	if err != nil || pageNum <= 0 {
		rootId = NotSet
	}
	authorId := query.Get("authorId")
	videoId := query.Get("videoId")

	return CommentsFilter{*pageSpec, rootId, authorId, videoId}, nil
}
