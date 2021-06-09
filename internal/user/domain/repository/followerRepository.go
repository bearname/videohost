package repository

type FollowerRepo interface {
	Follow(followingToUserId string, follower string, isFollowing bool) error
	GetStats(userId string) (int, error)
}
