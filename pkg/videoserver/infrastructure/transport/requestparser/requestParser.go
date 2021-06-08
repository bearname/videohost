package requestparser

import (
	"encoding/json"
	"errors"
	"github.com/bearname/videohost/pkg/common/util"
	dto2 "github.com/bearname/videohost/pkg/videoserver/domain/dto"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type RequestParser interface {
	Parse(r *http.Request) (interface{}, error)
}

type UploadVideoRequestParser struct {
}

func NewUploadVideoRequestParser() *UploadVideoRequestParser {
	return new(UploadVideoRequestParser)
}

func (p *UploadVideoRequestParser) Parse(request *http.Request) (interface{}, error) {
	title := request.FormValue("title")
	if len(title) == 0 {
		return nil, errors.New("cannot get title")
	}
	description := request.FormValue("description")
	if len(description) == 0 {
		return nil, errors.New("cannot get description")
	}

	fileReader, header, err := request.FormFile("file")
	if err != nil {
		return nil, errors.New("cannot get file")
	}
	chaptersJson := request.FormValue("chapters")
	var chapters []dto2.ChapterDto
	if len(chaptersJson) != 0 {
		if err = json.Unmarshal([]byte(chaptersJson), &chapters); err != nil {
			return nil, errors.New("cannot parse chapter")
		}
	}

	return &dto2.UploadVideoDto{
		Title:         title,
		Description:   description,
		MultipartFile: fileReader,
		FileHeader:    header,
		Chapters:      chapters,
	}, nil
}

type CatalogVideoParser struct {
}

func NewCatalogVideoParser() *CatalogVideoParser {
	return new(CatalogVideoParser)
}

func (p *CatalogVideoParser) Parse(request *http.Request) (interface{}, error) {
	var page int
	page, err := getIntQueryParameter(request, "page")
	if err != nil {
		return nil, errors.New("failed get page parameter")
	}

	var countVideoOnPage int
	countVideoOnPage, err = getIntQueryParameter(request, "countVideoOnPage")
	if err != nil {
		return nil, errors.New("failed get countVideoOnPage parameter")
	}

	return dto2.SearchDto{
		Page:         page,
		Count:        countVideoOnPage,
		SearchString: "",
	}, nil
}

type SearchVideoParser struct {
}

func NewSearchVideoParser() *SearchVideoParser {
	return new(SearchVideoParser)
}

func (p *SearchVideoParser) Parse(request *http.Request) (interface{}, error) {
	var page int
	page, err := getIntQueryParameter(request, "page")
	if err != nil {
		return nil, errors.New("failed get page parameter")
	}

	var countVideoOnPage int
	countVideoOnPage, err = getIntQueryParameter(request, "limit")
	if err != nil {
		return nil, errors.New("failed get countVideoOnPage parameter")
	}
	var searchString string
	searchString, ok := parseQueryParameter(request, "search")
	if !ok {
		return nil, errors.New("failed get searchString parameter")
	}

	return dto2.SearchDto{
		Page:         page,
		Count:        countVideoOnPage,
		SearchString: searchString,
	}, nil
}

func getIntQueryParameter(request *http.Request, key string) (int, error) {
	pageStr, done := parseQueryParameter(request, key)
	if !done {
		return 0, errors.New("invalid " + key + " parameter not found")
	}

	page, b := validate(pageStr)
	if b {
		return 0, errors.New("invalid " + key + " parameter not found")
	}

	return page, nil
}

func parseQueryParameter(request *http.Request, key string) (string, bool) {
	query := request.URL.Query()
	keys, ok := query[key]

	if !ok || len(keys) != 1 {
		return "", false
	}

	return keys[0], true
}

func validate(pageStr string) (int, bool) {
	page, b := util.StrToInt(pageStr)
	if !b || page < 0 {
		return 0, true
	}
	return page, false
}

type LikeVideoRequestParser struct {
}

func NewLikeVideoRequestParser() *LikeVideoRequestParser {
	return new(LikeVideoRequestParser)
}

func (p *LikeVideoRequestParser) Parse(request *http.Request) (interface{}, error) {
	authorId, ok := context.Get(request, "userId").(string)
	if !ok {
		return nil, errors.New("cannot get userId")
	}
	vars := mux.Vars(request)
	videoId, ok := vars["videoId"]
	if !ok {
		return nil, errors.New("videoId not present")
	}

	isLikeResult, ok := vars["isLike"]
	if !ok {
		return nil, errors.New("isLike not present")
	}

	atoi, err := strconv.Atoi(isLikeResult)
	if err != nil {
		return nil, errors.New("isLike must be 0 - is dislike or 1 - is like")
	}

	var isLike bool
	if atoi == 1 {
		isLike = true
	} else if atoi == 0 {
		isLike = false
	} else {
		return nil, errors.New("isLike must be 0 - is dislike or 1 - is like")
	}

	return &LikeVideoRequest{VideoId: videoId, OwnerId: authorId, IsLike: isLike}, nil
}
