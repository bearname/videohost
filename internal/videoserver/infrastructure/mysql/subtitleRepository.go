package mysql

import (
	"database/sql"
	"github.com/bearname/videohost/internal/common/db"
	"github.com/bearname/videohost/internal/videoserver/domain/dto"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
)

type SubtitleRepository struct {
	connector db.Connector
}

func NewSubtitleRepository(connector db.Connector) *SubtitleRepository {
	m := new(SubtitleRepository)
	m.connector = connector
	return m
}

func (r *SubtitleRepository) Create(subtitle dto.CreateSubtitleRequestDto) (int64, error) {
	var id int64
	err := db.WithTransaction(r.connector.GetDb(), func(tx db.Transaction) error {
		query := `INSERT INTO subtitle (video_id) VALUES (?);`

		var result sql.Result
		result, err := tx.Exec(`INSERT INTO subtitle (video_id) VALUES (?);`, subtitle.VideoId)
		if err != nil {
			return err
		}
		id, err = result.LastInsertId()
		if err != nil {
			return err
		}

		query = ""

		var values []interface{}
		for _, part := range subtitle.Parts {
			query += "INSERT INTO subtitle_part (start, end, text) VALUES (?, ?, ?);"
			values = append(values, part.Start, part.End, part.Text)
		}

		_, err = tx.Exec(query, values...)
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (r *SubtitleRepository) Find(videoId int) (model.Subtitle, error) {
	panic("implement subtitleRepo.findByVideo")
}

func (r *SubtitleRepository) Update(subtitle model.Subtitle) error {
	panic("implement subtitleRepo.Update")
}

func (r *SubtitleRepository) Delete(subtitleId int) error {
	panic("implement subtitleRepo.Delete")
}
