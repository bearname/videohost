package publisher

import (
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/thumbgenerator/domain/model"
	log "github.com/sirupsen/logrus"
)

func PublishTask(db database.Connector) *model.Task {
	var task model.Task
	err := db.GetDb().QueryRow("SELECT id_video, url FROM video WHERE status=?;", model.NotProcessed).Scan(
		&task.Id,
		&task.Url,
	)
	if err != nil {
		log.Info("Not found not processed task")
		log.Error(err.Error())
		return nil
	}
	log.Info("found not processed task " + task.Id)

	query := "UPDATE video SET status=? WHERE id_video=?;"
	err = db.ExecTransaction(query, model.Processing, task.Id)

	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return &task
}
