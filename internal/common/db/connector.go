package db

import (
	"database/sql"
)

type Connector interface {
	GetDb() *sql.DB
	Connect(user string, password string, dbAddress string, dbName string) error
	Close() error
	ExecTransaction(query string, args ...interface{}) error
}
