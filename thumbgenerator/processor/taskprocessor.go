package processor

import (
	"database/sql"
	"fmt"
	"github.com/bearname/videohost/thumbgenerator/model"
	"github.com/bearname/videohost/thumbgenerator/repository/mysql"
	"github.com/bearname/videohost/thumbgenerator/util"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ProcessTask(task *model.Task, db *sql.DB) {
	url := task.Url
	thumbUrl := filepath.Join(filepath.Dir(url), util.ThumbFileName)
	log.Info("ProcessTask" + task.Id + task.Url)
	fmt.Println("ProcessTask" + task.Id + task.Url)
	outputHls := url[0 : strings.LastIndex(url, "\\")+1]
	url = strings.ReplaceAll(url, "\\", "\\")
	root := "C:\\Users\\mikha\\go\\src\\videoserver\\videoserver\\"
	url = root + strings.ReplaceAll(url, "\\", "\\")
	outputHls = root + strings.ReplaceAll(outputHls, "\\", "\\")

	out, err := exec.Command("C:\\Users\\mikha\\go\\src\\videoserver\\videoprocessor\\videoprocessor.exe", url, root+strings.ReplaceAll(thumbUrl, "\\", "\\"), outputHls).Output()
	if err != nil {
		log.Error(err.Error())
		return
	}

	duration, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		log.Error(err.Error())
		return
	}

	//duration := rand.Intn(100)
	err = mysql.ExecTransaction(
		db,
		"UPDATE video SET status=?, duration=?, thumbnail_url=? WHERE id_video=?;",
		model.Processed,
		duration,
		thumbUrl,
		task.Id,
	)

	if err != nil {
		log.Error(err.Error())
	}
}
