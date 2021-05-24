package dto

type SignupUserDto struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	IsSubscribed bool   `json:"isSubscribed"`
}
