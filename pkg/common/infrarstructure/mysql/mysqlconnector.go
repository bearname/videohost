package mysql

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type ConnectorImpl struct {
	database *sql.DB
}

const user = "root"
const password = "123"
const databaseName = "video"

func (c *ConnectorImpl) GetDb() *sql.DB {
	return c.database
}

func (c *ConnectorImpl) Connect() error {
	if c.database != nil {
		log.Info("Already connected")
	}

	db, err := sql.Open("mysql", user+":"+password+"@/"+databaseName+"?parseTime=true")
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
