package mysql

import (
	"database/sql"
	"encoding/json"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/videoserver/app/dto"
	"github.com/bearname/videohost/pkg/videoserver/domain/model"
	log "github.com/sirupsen/logrus"
)

type VideoRepository struct {
	connector database.Connector
}

func NewMysqlVideoRepository(connector database.Connector) *VideoRepository {
	m := new(VideoRepository)
	m.connector = connector
	return m
}

func (r *VideoRepository) Create(userId string, videoId string, title string, description string, url string, chapters []dto.ChapterDto) error {
	if chapters == nil || len(chapters) == 0 {
		return r.insertWithoutChapter(userId, videoId, title, description, url)
	}
	return r.insertWithChapter(userId, videoId, title, description, url, chapters)
}

func (r *VideoRepository) insertWithChapter(userId string, videoId string, title string, description string, url string, chapters []dto.ChapterDto) error {
	var values []interface{}
	createChapterQuery := "INSERT INTO video_chapter (id_video, title, start, end) VALUES "

	for _, chapter := range chapters {
		createChapterQuery += "(?, ?, ?, ?), "
		values = append(values, videoId, chapter.Title, chapter.Start, chapter.End)
	}
	createChapterQuery = createChapterQuery[0 : len(createChapterQuery)-2]
	createChapterQuery += ";"

	err := WithTransaction(r.connector.GetDb(), func(tx Transaction) error {
		createVideoQuery := "INSERT INTO video (id_video, title, description, url, owner_id) VALUE (?, ?, ?, ?, ?); "
		_, err := tx.Exec(createVideoQuery, videoId, title, description, url, userId)
		if err != nil {
			return err
		}
		_, err = tx.Exec(createChapterQuery, values...)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type TxFn func(Transaction) error

func WithTransaction(db *sql.DB, fn TxFn) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
func (r *VideoRepository) insertWithoutChapter(userId string, videoId string, title string, description string, url string) error {
	query := "INSERT INTO video (id_video, title, description, url, owner_id) VALUE (?, ?, ?, ?, ?);"
	_, err := r.connector.GetDb().Query(query, videoId,
		title,
		description,
		url,
		userId)
	if err != nil {
		log.Info(err.Error())
		return err
	}
	return nil
}

func (r *VideoRepository) Save(video *model.Video) error {
	query := `INSERT INTO video (id_video, title, description, duration, status, thumbnail_url, url, uploaded, quality, owner_id) 
			VALUE (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err := r.connector.GetDb().Query(query, video.Id,
		video.Name,
		video.Description,
		video.Duration,
		video.Status,
		video.Thumbnail,
		video.Url,
		video.Uploaded,
		video.Quality,
		video.OwnerId)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (r *VideoRepository) Find(videoId string) (*model.Video, error) {
	var video model.Video
	q := `SELECT id_video,
       video.title AS  video_title,
       description AS video_description,
       duration AS video_duration,
       thumbnail_url AS video_thumbnail_url,
       url AS video_url,
       uploaded AS video_uploaded,
       quality AS video_quality,
       views AS video_views,
       owner_id AS video_owner_id,
       status AS video_status,
       GROUP_CONCAT(CONCAT('{"title":"', vc.title, '","start":', vc.start, ',"end":', vc.end, '}') SEPARATOR ',') AS video_chapters
FROM video
         LEFT JOIN video_chapter vc ON id_video = vc.video_id
WHERE id_video = ?
GROUP BY id_video;`
	row := r.connector.GetDb().QueryRow(q, videoId)
	var chapterString sql.NullString

	err := row.Scan(
		&video.Id,
		&video.Name,
		&video.Description,
		&video.Duration,
		&video.Thumbnail,
		&video.Url,
		&video.Uploaded,
		&video.Quality,
		&video.Views,
		&video.OwnerId,
		&video.Status,
		&chapterString,
	)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var chapters []model.Chapter

	if chapterString.Valid {
		chapters, err = r.parseChapter(chapterString.String, err)
		if err != nil {
			return nil, err
		}
	}

	video.Chapters = chapters

	return &video, err
}

func (r *VideoRepository) setChapter(chapterString string, err error, video model.Video) (model.Video, error) {
	var chapters []model.Chapter

	if len(chapterString) > 0 {
		chapters, err = r.parseChapter(chapterString, err)
		if err != nil {
			return video, err
		}

	}

	video.Chapters = chapters
	return video, nil
}

func (r *VideoRepository) parseChapter(chapterString string, err error) ([]model.Chapter, error) {
	var chapters []model.Chapter

	chapterString = "[" + chapterString + "]"
	err = json.Unmarshal([]byte(chapterString), &chapters)
	return chapters, err
}

func (r *VideoRepository) Update(videoId string, title string, description string) error {
	rows, err := r.connector.GetDb().Query("UPDATE video SET title=?, description=? WHERE id_video=?;", title, description, videoId)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer rows.Close()
	return nil
}

func (r *VideoRepository) Delete(videoId string) error {
	rows, err := r.connector.GetDb().Query("DELETE FROM video  WHERE id_video=?;", videoId)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (r *VideoRepository) FindVideosByPage(page int, count int) ([]model.VideoListItem, error) {
	offset := (page) * count
	query := "SELECT id_video, title, duration, thumbnail_url, uploaded, views, status, quality FROM video WHERE status=3 LIMIT ?, ?;"
	rows, err := r.connector.GetDb().Query(query, offset, count)

	return r.getVideoListItem(rows, err)
}

func (r *VideoRepository) GetPageCount(countVideoOnPage int) (int, bool) {
	query := "SELECT COUNT(id_video) AS countReadyVideo FROM video WHERE status=3;"
	rows, err := r.connector.GetDb().Query(query)
	if err != nil {
		return 0, false
	}
	defer rows.Close()

	var countVideo int
	for rows.Next() {
		err := rows.Scan(
			&countVideo,
		)
		if err != nil {
			return 0, false
		}
	}
	countPage := countVideo / countVideoOnPage
	if countVideo%countVideoOnPage > 0 {
		countPage += 1
	}
	return countPage, true
}

func (r *VideoRepository) AddVideoQuality(videoId string, quality string) error {
	query := "UPDATE video SET `quality` = concat(quality,  concat(',',  ?))  WHERE id_video = ?;"
	rows, err := r.connector.GetDb().Query(query, quality, videoId)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (r *VideoRepository) SearchVideo(searchString string, page int, count int) ([]model.VideoListItem, error) {

	offset := (page - 1) * count
	query := `SELECT id_video, title, duration, thumbnail_url, uploaded, views, status, quality  FROM video 
               WHERE MATCH(video.title) AGAINST (? IN NATURAL LANGUAGE MODE) AND status=3 LIMIT ?, ?;`
	rows, err := r.connector.GetDb().Query(query, searchString, offset, count)

	return r.getVideoListItem(rows, err)
}

func (r *VideoRepository) IncrementViews(id string) bool {
	query := "UPDATE video SET video.views = video.views + 1 WHERE id_video=?"
	rows, err := r.connector.GetDb().Query(query, id)
	if err != nil {
		log.Info(err.Error())
		return false
	}
	defer rows.Close()
	return true
}

func (r *VideoRepository) FindUserVideos(userId string, dto dto.SearchDto) ([]model.VideoListItem, error) {
	offset := (dto.Page) * dto.Count
	query := "SELECT video.id_video, title, duration, thumbnail_url, uploaded, views, status, quality FROM video  WHERE owner_id=?  LIMIT ?, ?;"
	rows, err := r.connector.GetDb().Query(query, userId, offset, dto.Count)
	return r.getVideoListItem(rows, err)
}

func (r *VideoRepository) getVideoListItem(rows *sql.Rows, err error) ([]model.VideoListItem, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := make([]model.VideoListItem, 0)
	for rows.Next() {
		var videoListItem model.VideoListItem
		err := rows.Scan(
			&videoListItem.Id,
			&videoListItem.Name,
			&videoListItem.Duration,
			&videoListItem.Thumbnail,
			&videoListItem.Uploaded,
			&videoListItem.Views,
			&videoListItem.Status,
			&videoListItem.Quality,
		)
		if err != nil {
			return nil, err
		}
		videos = append(videos, videoListItem)
	}

	return videos, nil
}
