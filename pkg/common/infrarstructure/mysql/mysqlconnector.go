package mysql

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type ConnectorImpl struct {
	database *sql.DB
}

func NewConnectorImpl() ConnectorImpl {
	return *new(ConnectorImpl)
}

func (c *ConnectorImpl) GetDb() *sql.DB {
	return c.database
}

func (c *ConnectorImpl) Connect(user string, password string, dbAddress string, dbName string) error {
	if c.database != nil {
		log.Info("Already connected")
	}

	dataSourceName := user + ":" + password + "@"
	if len(dbAddress) != 0 {
		dataSourceName += dbAddress
	}

	dataSourceName += "/" + dbName + "?parseTime=true"

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Error(err)
		return err
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return err
	}

	c.database = db

	return nil
}

func (c *ConnectorImpl) Close() error {
	err := c.database.Close()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	c.database = nil

	return nil
}

func (c *ConnectorImpl) ExecTransaction(query string, args ...interface{}) error {
	tx, err := c.database.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Error(err)
		}
	}(tx)

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
