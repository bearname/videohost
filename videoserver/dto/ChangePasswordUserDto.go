package dto

type ChangePasswordUserDto struct {
	Username    string `json:"username"`
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}
