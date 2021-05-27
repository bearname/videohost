package dto

type FindUserDto struct {
	Email        string `json:"email"`
	Username     string `json:"usernameOrId"`
	IsSubscribed bool   `json:"isSubscribed"`
	Role         int    `json:"role"`
}
