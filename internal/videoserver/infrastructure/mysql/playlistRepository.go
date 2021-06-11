package mysql

import (
	"database/sql"
	"github.com/bearname/videohost/internal/common/db"
	"github.com/bearname/videohost/internal/videoserver/domain"
	"github.com/bearname/videohost/internal/videoserver/domain/dto"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
	"strconv"
)

type PlaylistRepository struct {
	connector db.Connector
}

func NewPlaylistRepository(connector db.Connector) *PlaylistRepository {
	m := new(PlaylistRepository)
	m.connector = connector
	return m
}

func (r *PlaylistRepository) Create(playlist dto.CreatePlaylistDto) (int64, error) {
	_, err := r.checkExist(playlist.OwnerId, playlist.Name)
	if err != nil && err != domain.ErrPlaylistNotFound {
		return 0, err
	}

	var id int64
	err = WithTransaction(r.connector.GetDb(), func(tx Transaction) error {
		query := `INSERT INTO playlist (name, user_id, privacy) VALUES (?, ?, ?);`
		var result sql.Result
		result, err = tx.Exec(query, playlist.Name, playlist.OwnerId, playlist.Privacy)
		if err != nil {
			return err
		}
		id, err = result.LastInsertId()
		if err != nil {
			return err
		}
		queryVideos := "INSERT INTO video_in_playlist (playlist_id, video_id, user_id) VALUES "
		for i, videoId := range playlist.Videos {
			queryVideos += "(" + strconv.FormatInt(id, 10) + ",'" + videoId + "','" + playlist.OwnerId + "')"
			if i != len(playlist.Videos)-1 {
				queryVideos += ","
			} else {
				queryVideos += ";"
			}
		}

		_, err = tx.Exec(queryVideos)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PlaylistRepository) checkExist(ownerId string, name string) (int, error) {
	sqlQuery := `SELECT id FROM playlist WHERE user_id = ? AND name = ?`
	rows, err := r.connector.GetDb().Query(sqlQuery, ownerId, name)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var id int
	if rows.Next() {
		rows.Scan(&id)
		return id, nil
	}
	return 0, domain.ErrPlaylistNotFound
}

func (r *PlaylistRepository) FindPlaylist(playlistId int) (model.Playlist, error) {
	sqlQuery := `SELECT p.id,
       				p.user_id,
				    p.name,
				    p.created,
				    p.privacy,
					   GROUP_CONCAT(CONCAT('{"videoId":"', video_in_playlist.video_id, '}') SEPARATOR ',') AS video_chapters
				FROM video_in_playlist
						 LEFT JOIN playlist p on video_in_playlist.playlist_id = p.id
				WHERE playlist_id = ?
				GROUP BY playlist_id;`

	rows, err := r.connector.GetDb().Query(sqlQuery, playlistId)
	if err != nil {
		return model.Playlist{}, err
	}
	defer rows.Close()
	var playlist model.Playlist

	if rows.Next() {
		err = rows.Scan(&playlist.Id,
			&playlist.OwnerId,
			&playlist.Name,
			&playlist.Created,
			&playlist.Privacy,
			&playlist.VideoString)

		if err != nil {
			return model.Playlist{}, err
		}

		return playlist, nil
	}
	return playlist, domain.ErrPlaylistNotFound
}

func (r *PlaylistRepository) AddVideos(playlistId int, userId string, videosId []string) error {
	query := "INSERT INTO video_in_playlist (playlist_id, video_id, user_id) VALUES "
	var vals []interface{}
	for i, videoId := range videosId {
		query += "(?, ?, ?)"
		vals = append(vals, playlistId, videoId, userId)

		if i == len(videosId)-1 {
			query += ";"
		} else {
			query += ","
		}
	}

	_, err := r.connector.GetDb().Query(query, vals...)
	if err != nil {
		return err
	}
	return nil
}

func (r *PlaylistRepository) RemoveVideos(playlistId int, userId string, videosId []string) error {
	query := "DELETE FROM video_in_playlist WHERE playlist_id = ? AND user_id = ? AND video_id  IN ("
	length := len(videosId) - 1
	var vals []interface{}
	vals = append(vals, playlistId, userId)

	for i, videoId := range videosId {
		query += "?"
		vals = append(vals, videoId)
		if i != length {
			query += ","
		}
	}
	query += ");"

	_, err := r.connector.GetDb().Query(query, vals...)
	if err != nil {
		return err
	}
	return nil
}

func (r *PlaylistRepository) ChangeOrder(playlistId int, videoId []string, order []int) error {

	return nil
}

func (r *PlaylistRepository) ChangePrivacy(ownerId string, playlistId int, privacyType model.PrivacyType) error {
	query := "UPDATE playlist SET privacy = ?  WHERE id = ? AND user_id=?;"
	err := r.connector.ExecTransaction(query, privacyType, playlistId, ownerId)

	return err
}

func (r *PlaylistRepository) Delete(ownerId string, playlistId int) error {
	rows, err := r.connector.GetDb().Query("DELETE FROM playlist WHERE id=? AND user_id=?;", playlistId, ownerId)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
