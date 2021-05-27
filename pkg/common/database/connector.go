package database

import "database/sql"

type Connector interface {
	GetDb() *sql.DB
	Connect() error
	Close() error
	ExecTransaction(query string, args ...interface{}) error
}
