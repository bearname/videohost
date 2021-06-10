package mysql

import (
	"database/sql"
	"errors"
	"github.com/bearname/videohost/internal/common/db"
	"github.com/bearname/videohost/internal/video-comments/domain"
	log "github.com/sirupsen/logrus"
)

type CommentRepo struct {
	connector db.Connector
}

func NewCommentRepo(connector db.Connector) *CommentRepo {
	m := new(CommentRepo)
	m.connector = connector
	return m
}

func (r *CommentRepo) Create(VideoId string, UserId string, ParentId int, Message string) (int64, error) {
	var sqlQuery string
	var err error
	var result sql.Result
	if ParentId > 0 {
		sqlQuery = "INSERT INTO video_comments (video_id, user_id, parent_id, message) VALUES  (?, ?, ?, ?);"
		result, err = r.connector.GetDb().Exec(sqlQuery, &VideoId, &UserId, &ParentId, &Message)
	} else {
		sqlQuery = "INSERT INTO video_comments (video_id, user_id, message) VALUES  (?, ?, ?);"
		result, err = r.connector.GetDb().Exec(sqlQuery, &VideoId, &UserId, &Message)
	}

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *CommentRepo) FindById(commentId int) (domain.Comment, error) {
	query := `SELECT video_comments.id, video_id, user_id, message, created FROM video_comments WHERE id=?;`

	rows, err := r.connector.GetDb().Query(query, commentId)
	if err != nil {
		return domain.Comment{}, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Error(err)
		}
	}(rows)

	if rows.Next() {
		var comment domain.Comment
		err = rows.Scan(
			&comment.Id,
			&comment.VideoId,
			&comment.UserId,
			&comment.Message,
			&comment.Created,
		)

		if err != nil {
			return domain.Comment{}, err
		}

		return comment, nil
	}
	return domain.Comment{}, errors.New("not found comment")
}

func (r *CommentRepo) FindRootLevel(videoId string, page *db.Page) (domain.VideoComments, error) {
	query := `SELECT id,  user_id,  message, created, countSubComments
				FROM video_comments
						 JOIN (WITH RECURSIVE tmp (id, parent_id) AS (SELECT id, parent_id
																	  FROM video_comments
																	  WHERE parent_id IS NULL
																		AND video_id = ?
																	  UNION ALL
																	  SELECT l.id, l.parent_id
																	  FROM tmp AS p
																			   JOIN video_comments AS l ON p.id = l.parent_id
)
               SELECT parent_id, COUNT(id) AS countSubComments
               FROM tmp
               GROUP BY parent_id
               HAVING countSubComments
               ) AS tmp ON tmp.parent_id = video_comments.parent_id
				LIMIT ?, ?;`

	rows, err := r.connector.GetDb().Query(query, videoId, page.Number, page.Size)
	if err != nil {
		return domain.VideoComments{}, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Error(err)
		}
	}(rows)

	comments := make([]domain.RootComment, 0)
	for rows.Next() {
		var comment domain.RootComment
		err = rows.Scan(
			&comment.Id,
			&comment.UserId,
			&comment.Message,
			&comment.Created,
			&comment.CountSubComments)
		if err != nil {
			return domain.VideoComments{}, err
		}
		comments = append(comments, comment)
	}

	rows, err = r.connector.GetDb().Query("SELECT COUNT(id) FROM video_comments;")
	if err != nil {
		return domain.VideoComments{}, err
	}
	defer rows.Close()
	var countAllComments int
	if rows.Next() {
		err = rows.Scan(&countAllComments)
		if err != nil {
			return domain.VideoComments{}, err
		}
	}

	return *domain.NewVideoComments(videoId, countAllComments, comments), nil
}

func (r *CommentRepo) FindUserComments(userId string, page *db.Page) (domain.Comments, error) {
	query := `SELECT id,  video_id,  message, created, countSubComments
FROM video_comments
         JOIN (WITH RECURSIVE tmp (id, parent_id) AS (SELECT id, parent_id
                                                      FROM video_comments
                                                      WHERE parent_id IS NULL
                                                        AND user_id = ?
                                                      UNION ALL
                                                      SELECT l.id, l.parent_id
                                                      FROM tmp AS p
                                                               JOIN video_comments AS l ON p.id = l.parent_id
)
               SELECT parent_id, COUNT(id) AS countSubComments
               FROM tmp
               GROUP BY parent_id
               HAVING countSubComments
) AS tmp ON tmp.parent_id = video_comments.parent_id
LIMIT ?, ?;`

	rows, err := r.connector.GetDb().Query(query, userId, page.Number, page.Size)
	if err != nil {
		return domain.Comments{}, err
	}
	defer rows.Close()

	comments := make([]domain.RootComment, 0)
	for rows.Next() {
		var comment domain.RootComment
		err = rows.Scan(
			&comment.Id,
			&comment.VideoId,
			&comment.Message,
			&comment.Created,
			&comment.CountSubComments)
		if err != nil {
			return domain.Comments{}, err
		}
		comments = append(comments, comment)
	}

	rows, err = r.connector.GetDb().Query("SELECT COUNT(id) FROM video_comments WHERE user_id=?;", userId)
	if err != nil {
		return domain.Comments{}, err
	}
	defer rows.Close()
	var countAllComments int
	if rows.Next() {
		err = rows.Scan(&countAllComments)
		if err != nil {
			return domain.Comments{}, err
		}
	}

	return domain.Comments{CountAllComments: countAllComments, RootComments: comments}, nil
}

func (r *CommentRepo) FindChildren(rootCommentId int, page *db.Page) ([]domain.Comment, error) {
	query := `SELECT video_comments.id, video_id, user_id, message, created
				FROM video_comments
						 JOIN (
					WITH RECURSIVE tmp (id) AS (SELECT id
												FROM video_comments
												WHERE parent_id IS NULL AND parent_id = ?
												UNION ALL
												SELECT l.id
												FROM tmp AS p
														 JOIN video_comments AS l
															  ON p.id = l.parent_id
					)
					SELECT id
					FROM tmp
					LIMIT ?, ?
				) AS ids
							  ON ids.id = video_comments.id;`

	rows, err := r.connector.GetDb().Query(query, rootCommentId, page.Number, page.Size)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Error(err)
		}
	}(rows)

	comments := make([]domain.Comment, 0)
	for rows.Next() {
		var comment domain.Comment
		err = rows.Scan(
			&comment.Id,
			&comment.VideoId,
			&comment.UserId,
			&comment.Message,
			&comment.Created,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepo) Edit(commentId int, message string) error {
	rows, err := r.connector.GetDb().Query("UPDATE video_comments SET message=?, created=NOW() WHERE id=?;", message, commentId)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer rows.Close()
	return nil
}

func (r *CommentRepo) Delete(commentId int) error {
	rows, err := r.connector.GetDb().Query("DELETE FROM video_comments  WHERE id=?;", commentId)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
