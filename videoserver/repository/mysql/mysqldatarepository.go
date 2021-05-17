package mysql

import (
	"github.com/bearname/videohost/videoserver/model"
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

	row := r.connector.Database.QueryRow("SELECT id_video, title, description, duration, thumbnail_url, url, uploaded FROM video WHERE id_video=? ORDER BY uploaded DESC", id)
	//i := len("C:\\Users\\mikha\\go\\src\\videoserver\\videoserver\\content\\")

	err := row.Scan(
		&video.Id,
		&video.Name,
		&video.Description,
		&video.Duration,
		&video.Thumbnail,
		&video.Url,
		&video.Uploaded,
	)
	//thumbnail := video.Thumbnail[i:len(video.Thumbnail)]
	//video.Thumbnail = thumbnail
	return &video, err
}

func (r *DataRepository) GetVideoList(page int, count int) ([]model.VideoListItem, error) {
	args := (page) * count
	rows, err := r.connector.Database.Query("SELECT id_video, title, duration, thumbnail_url FROM video WHERE status=3 LIMIT ?, ?;", args, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := make([]model.VideoListItem, 0)
	for rows.Next() {
		var videoListItem model.VideoListItem
		//i := len("C:\\Users\\mikha\\go\\src\\videoserver\\videoserver\\content\\")
		err := rows.Scan(
			&videoListItem.Id,
			&videoListItem.Name,
			&videoListItem.Duration,
			&videoListItem.Thumbnail,
		)
		//thumbnail := videoListItem.Thumbnail[i:len(videoListItem.Thumbnail)]
		//videoListItem.Thumbnail = "http://localhost:8000/content/" + thumbnail
		if err != nil {
			return nil, err
		}
		videos = append(videos, videoListItem)
	}

	return videos, nil
}

func (r *DataRepository) NewVideo(id string, fileName string, description string, url string) error {
	return ExecTransaction(
		r.connector.Database,
		"INSERT INTO video SET id_video=?, title=?, description=?, url=?;",
		id,
		fileName,
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

func (r *DataRepository) AddVideoQuality(id string, quality string) {
	rows, err := r.connector.Database.Query("UPDATE video SET `quality` = concat(quality, ?)  WHERE id = ?;", quality, id)
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
