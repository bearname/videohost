package mysql

import (
	"github.com/bearname/videohost/internal/common/db"
)

type FollowerRepo struct {
	connector db.Connector
}

func NewFollowerRepo(connector db.Connector) *FollowerRepo {
	m := new(FollowerRepo)
	m.connector = connector
	return m
}

func (r *FollowerRepo) Follow(followingToUserId string, followerId string, isFollowing bool) error {
	sqlQuery := `INSERT INTO subscriptions (following_to_user_id, follower_user_id, is_following) VALUES (?, ?, ?) 
		ON DUPLICATE KEY  UPDATE is_following=?;`
	query, err := r.connector.GetDb().Query(sqlQuery, followingToUserId, followerId, isFollowing, isFollowing)
	if err != nil {
		return err
	}

	defer query.Close()
	return nil
}

func (r *FollowerRepo) GetStats(userId string) (int, error) {
	query := "SELECT COUNT(follower_user_id) AS countFollowers FROM subscriptions WHERE following_to_user_id=?;"
	rows, err := r.connector.GetDb().Query(query, userId)

	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var countFollowers int

	for rows.Next() {
		err = rows.Scan(&countFollowers)
		if err != nil {
			return 0, err
		}
	}

	return countFollowers, nil
}
