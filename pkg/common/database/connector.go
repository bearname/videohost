package database

type Connector interface {
	Connect() error
	Close() error
}
