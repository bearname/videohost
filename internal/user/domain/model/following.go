package model

type Following struct {
	FollowingToUserId string
	FollowerUserId    string
}

type Subscription struct {
	UserId   string
	UserName string
}

type UserStatistic struct {
	Subscription
	CountSubscription int
}

func NewUserStatistic(userId string, userName string, countSubscription int) *UserStatistic {
	u := new(UserStatistic)
	u.UserId = userId
	u.UserName = userName
	u.CountSubscription = countSubscription
	return u
}
