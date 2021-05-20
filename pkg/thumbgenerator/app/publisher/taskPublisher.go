package publisher

import (
	"database/sql"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/thumbgenerator/domain/model"
	log "github.com/sirupsen/logrus"
)

func PublishTask(db *sql.DB) *model.Task {
	var task model.Task
	err := db.QueryRow("SELECT id_video, url FROM video WHERE status=?;", model.NotProcessed).Scan(
		&task.Id,
		&task.Url,
	)
	if err != nil {
		log.Info("Not found not processed task")
		log.Error(err.Error())
		return nil
	}
	log.Info("found not processed task " + task.Id)

	err = database.ExecTransaction(
		db,
		"UPDATE video SET status=? WHERE id_video=?;", model.Processing,
		task.Id,
	)

	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return &task
}
