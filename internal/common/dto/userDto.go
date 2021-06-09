package dto

type UserDto struct {
	Username string `json:"username"`
	UserId   string `json:"user_id"`
	Ok       bool   `json:"ok"`
	Role     int    `json:"role"`
	Token    string `json:"token"`
}
