package database

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type Connector struct {
	Database *sql.DB
}

const user = "root"
const password = "123"
const databaseName = "video"

func (c *Connector) Connect() error {
	if c.Database != nil {
		log.Info("Already connected")
	}

	db, err := sql.Open("mysql", user+":"+password+"@/"+databaseName)
	if err != nil {
		log.Error(err)
		return err
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return err
	}

	c.Database = db

	return nil
}

func (c *Connector) Close() error {
	err := c.Database.Close()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	c.Database = nil

	return nil
}

func ExecTransaction(db *sql.DB, query string, args ...interface{}) error {
	tx, err := db.Begin()
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
