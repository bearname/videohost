package repository

type UserRepository struct {
	Users    map[string][]byte
	idxUsers int
}

func NewUserRepository() *UserRepository {
	v := new(UserRepository)

	v.Users = make(map[string][]byte)
	v.idxUsers = 0

	return v
}
