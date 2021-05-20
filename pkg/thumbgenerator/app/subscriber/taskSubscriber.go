package subscriber

import (
	"database/sql"
	"fmt"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/common/util"
	"github.com/bearname/videohost/pkg/thumbgenerator/domain/model"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func HandleTask(task *model.Task, db *sql.DB) {
	url := task.Url
	thumbUrl := filepath.Join(filepath.Dir(url), util.ThumbFileName)
	log.Info("HandleTask" + task.Id + task.Url)
	fmt.Println("HandleTask" + task.Id + task.Url)
	outputHls := url[0 : strings.LastIndex(url, "\\")+1]
	url = strings.ReplaceAll(url, "\\", "\\")
	root := "C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\"
	url = root + strings.ReplaceAll(url, "\\", "\\")
	outputHls = root + strings.ReplaceAll(outputHls, "\\", "\\")
	out, err := exec.Command("C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoprocessor\\videoprocessor.exe", url, root+strings.ReplaceAll(thumbUrl, "\\", "\\"), outputHls).Output()
	if err != nil {
		log.Error(err.Error())
		return
	}

	duration, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = database.ExecTransaction(
		db,
		"UPDATE video SET status=?, duration=?, thumbnail_url=? WHERE id_video=?;",
		model.Processed,
		duration,
		thumbUrl,
		task.Id,
	)

	if err != nil {
		log.Error("Failed set status processed" + err.Error())
	}
}
