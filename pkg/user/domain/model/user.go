package model

type Role int

const (
	Admin Role = iota
	General
)

func (q Role) String() string {
	return [...]string{"Admin", "General"}[q]
}

func (q Role) Values() int {
	return [...]int{0, 1}[q]
}

type User struct {
	Id           int    `json:"id"`
	Key          string `json:"key"`
	Username     string `json:"username"`
	Password     []byte `json:"password"`
	Email        string `json:"email"`
	IsSubscribed bool   `json:"isSubscribed"`
	Created      string `json:"uploaded"`
	Role         Role   `json:"role"`
	Secret       string `json:"secret"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
