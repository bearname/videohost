package mysql

import (
	"database/sql"
	"github.com/bearname/videohost/videoserver/model"
	log "github.com/sirupsen/logrus"
)

type DataRepository struct {
	connector Connector
}

func NewMysqlDataRepository(connector Connector) *DataRepository {
	m := new(DataRepository)
	m.connector = connector
	return m
}

func (r *DataRepository) GetVideo(id string) (*model.Video, error) {
	var video model.Video

	row := r.connector.Database.QueryRow("SELECT id_video, title, description, duration, thumbnail_url, url, uploaded, quality, views FROM video WHERE id_video=? AND status = 3 ORDER BY uploaded DESC", id)

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
	)
	r.IncrementViews(id)
	return &video, err
}

func (r *DataRepository) GetVideoList(page int, count int) ([]model.VideoListItem, error) {
	offset := (page) * count
	rows, err := r.connector.Database.Query("SELECT id_video, title, duration, thumbnail_url, uploaded, views FROM video WHERE status=3 LIMIT ?, ?;", offset, count)

	return r.getVideoListItem(rows, err)
}

func (r *DataRepository) NewVideo(id string, title string, description string, url string) error {
	return ExecTransaction(
		r.connector.Database,
		"INSERT INTO video SET id_video=?, title=?, description=?, url=?;",
		id,
		title,
		description,
		url,
	)
}

func (r *DataRepository) GetPageCount(countVideoOnPage int) (int, bool) {
	rows, err := r.connector.Database.Query("SELECT COUNT(id_video) AS countReadyVideo FROM video WHERE status=3;")
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

func (r *DataRepository) AddVideoQuality(id string, quality string) bool {
	rows, err := r.connector.Database.Query("UPDATE video SET `quality` = concat(quality,  concat(',',  ?))  WHERE id_video = ?;", quality, id)
	if err != nil {
		log.Info(err.Error())
		return false
	}
	defer rows.Close()
	return true
}

func (r *DataRepository) SearchVideo(searchString string, page int, count int) ([]model.VideoListItem, error) {

	offset := (page - 1) * count
	rows, err := r.connector.Database.Query("SELECT id_video, title, duration, thumbnail_url, uploaded, views FROM video WHERE MATCH(video.title) AGAINST (? IN NATURAL LANGUAGE MODE) AND status=3 LIMIT ?, ?;", searchString, offset, count)

	return r.getVideoListItem(rows, err)
}

func (r *DataRepository) IncrementViews(id string) bool {
	rows, err := r.connector.Database.Query("UPDATE video SET video.views = video.views + 1 WHERE id_video=?", id)
	if err != nil {
		log.Info(err.Error())
		return false
	}
	defer rows.Close()
	return true
}

func (r *DataRepository) getVideoListItem(rows *sql.Rows, err error) ([]model.VideoListItem, error) {

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
		)
		if err != nil {
			return nil, err
		}
		videos = append(videos, videoListItem)
	}

	return videos, nil
}
