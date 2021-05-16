package processor

import (
	"database/sql"
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
	log.Info("generete thumnail")
	outputHls := url[0 : strings.LastIndex(url, "\\")+1]
	url = strings.ReplaceAll(url, "\\", "\\")
	root := "C:\\Users\\mikha\\go\\src\\videoserver\\videoserver\\"
	url = root + strings.ReplaceAll(url, "\\", "\\")
	thumbUrl = root + strings.ReplaceAll(thumbUrl, "\\", "\\")
	outputHls = root + strings.ReplaceAll(outputHls, "\\", "\\")
	out, err := exec.Command("C:\\Users\\mikha\\go\\src\\videoserver\\videoprocessor\\videoprocessor.exe", url, thumbUrl, outputHls).Output()
	//out, err := exec.Command("C:\\Users\\mikha\\go\\src\\videoserver\\videoprocessor\\videoprocessor.exe",
	//	"C:\\Users\\mikha\\go\\src\\videoserver\\videoserver\\content\\a55a4477-b446-11eb-875e-00ff7c2a75d7\\index.mp4",
	//	"C:\\Users\\mikha\\go\\src\\videoserver\\videoserver\\content\\a55a4477-b446-11eb-875e-00ff7c2a75d7\\default.jpg",
	//	"C:\\Users\\mikha\\go\\src\\videoserver\\videoserver\\content\\a55a4477-b446-11eb-875e-00ff7c2a75d7\\").Output()
	if err != nil {
		log.Error(err.Error())
		return
	}
	duration, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		log.Error(err.Error())
		return
	}

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
