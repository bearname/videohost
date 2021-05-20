package mysql

import (
	"database/sql"
	"github.com/bearname/videohost/pkg/common/database"
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

func (r *VideoRepository) GetVideo(id string) (*model.Video, error) {
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

func (r *VideoRepository) GetVideoList(page int, count int) ([]model.VideoListItem, error) {
	offset := (page) * count
	rows, err := r.connector.Database.Query("SELECT id_video, title, duration, thumbnail_url, uploaded, views, status FROM video WHERE status=3 LIMIT ?, ?;", offset, count)

	return r.getVideoListItem(rows, err)
}

func (r *VideoRepository) NewVideo(userId string, videoId string, title string, description string, url string) error {
	_, err := r.connector.Database.Query("INSERT INTO video (id_video, title, description, url) VALUE (?, ?, ?, ?);", videoId,
		title,
		description,
		url)
	if err != nil {
		log.Info(err.Error())
		return err
	}
	_, err = r.connector.Database.Query("INSERT INTO user_videos (key_user, id_video) VALUE (?, ?);", userId, videoId)
	if err != nil {
		log.Info(err.Error())
		return err
	}

	return nil
}

func (r *VideoRepository) GetPageCount(countVideoOnPage int) (int, bool) {
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

func (r *VideoRepository) AddVideoQuality(id string, quality string) bool {
	rows, err := r.connector.Database.Query("UPDATE video SET `quality` = concat(quality,  concat(',',  ?))  WHERE id_video = ?;", quality, id)
	if err != nil {
		log.Info(err.Error())
		return false
	}
	defer rows.Close()
	return true
}

func (r *VideoRepository) SearchVideo(searchString string, page int, count int) ([]model.VideoListItem, error) {

	offset := (page - 1) * count
	rows, err := r.connector.Database.Query("SELECT id_video, title, duration, thumbnail_url, uploaded, views, status FROM video WHERE MATCH(video.title) AGAINST (? IN NATURAL LANGUAGE MODE) AND status=3 LIMIT ?, ?;", searchString, offset, count)

	return r.getVideoListItem(rows, err)
}

func (r *VideoRepository) IncrementViews(id string) bool {
	rows, err := r.connector.Database.Query("UPDATE video SET video.views = video.views + 1 WHERE id_video=?", id)
	if err != nil {
		log.Info(err.Error())
		return false
	}
	defer rows.Close()
	return true
}

func (r *VideoRepository) FindUserVideos(userId string, page int, count int) ([]model.VideoListItem, error) {
	offset := (page) * count
	query := "SELECT video.id_video, title, duration, thumbnail_url, uploaded, views, status FROM video LEFT JOIN user_videos uv on video.id_video = uv.id_video WHERE uv.key_user=?  LIMIT ?, ?;"
	rows, err := r.connector.Database.Query(query, userId, offset, count)
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
		)
		if err != nil {
			return nil, err
		}
		videos = append(videos, videoListItem)
	}

	return videos, nil
}
