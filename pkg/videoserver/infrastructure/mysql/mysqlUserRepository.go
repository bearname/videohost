package mysql

import (
	"github.com/bearname/videohost/pkg/common/database"
	model2 "github.com/bearname/videohost/pkg/videoserver/domain/model"
)

type UserRepository struct {
	connector database.Connector
}

func NewMysqlUserRepository(connector database.Connector) *UserRepository {
	m := new(UserRepository)
	m.connector = connector
	return m
}

func (r *UserRepository) CreateUser(key string, username string, password []byte, role model2.Role, accessToken string, refreshToken string) error {
	query, err := r.connector.Database.Query("INSERT INTO users (key_user, username, password, role, access_token, refresh_token) VALUES (?, ?, ?, ?, ?, ?);", key, username, password, role, accessToken, refreshToken)
	if err != nil {
		return err
	}

	defer query.Close()

	return nil
}

func (r *UserRepository) FindByUserName(username string) (*model2.User, error) {
	var user model2.User

	row := r.connector.Database.QueryRow("SELECT key_user, username, password, created, role, secret, access_token, refresh_token FROM users WHERE username = ?;", username)

	err := row.Scan(
		&user.Key,
		&user.Username,
		&user.Password,
		&user.Created,
		&user.Role,
		&user.Secret,
		&user.AccessToken,
		&user.RefreshToken,
	)

	return &user, err
}

func (r *UserRepository) UpdatePassword(username string, password []byte) bool {
	err := database.ExecTransaction(
		r.connector.Database,
		"UPDATE users SET password = ? WHERE username = ?;", password, username)

	return err == nil
}

func (r *UserRepository) UpdateAccessToken(username string, token string) bool {
	err := database.ExecTransaction(
		r.connector.Database,
		"UPDATE users SET access_token = ?  WHERE username = ?;", token, username)
	return err == nil
}

func (r *UserRepository) UpdateRefreshToken(username string, token string) bool {
	err := database.ExecTransaction(
		r.connector.Database,
		"UPDATE users SET refresh_token = ?  WHERE username = ?;", token, username)
	return err == nil
}

func (r *UserRepository) GetCountVideos(userId string) (int, bool) {
	rows, err := r.connector.Database.Query("SELECT COUNT(key_user) AS countReadyVideo FROM user_videos WHERE key_user=?;", userId)
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

	return countVideo, true
}
