package model

type DSN struct {
	Address  string
	Username string
	Password string
	DB       string
}

func NewDsn(address string, username string, password string, db string) DSN {
	return DSN{
		address, username, password, db,
	}
}
