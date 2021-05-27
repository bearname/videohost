package subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/bearname/videohost/pkg/common/caching"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/common/util"
	"github.com/bearname/videohost/pkg/thumbgenerator/domain/model"
	videoModel "github.com/bearname/videohost/pkg/videoserver/domain/model"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func HandleTask(task *model.Task, connector database.Connector, cache caching.Cache) {
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

	query := "UPDATE video SET status=?, duration=?, thumbnail_url=? WHERE id_video=?;"
	err = connector.ExecTransaction(
		query,
		model.Processed,
		duration,
		thumbUrl,
		task.Id,
	)

	var video videoModel.Video
	query = "SELECT id_video, title, description, duration, thumbnail_url, url, uploaded, quality, views, owner_id, status FROM video WHERE id_video=?;"
	row := connector.GetDb().QueryRow(query, task.Id)

	err = row.Scan(
		&video.Id,
		&video.Name,
		&video.Description,
		&video.Duration,
		&video.Thumbnail,
		&video.Url,
		&video.Uploaded,
		&video.Quality,
		&video.Views,
		&video.OwnerId,
		&video.Status,
	)
	if err == nil {
		marshal, err := json.Marshal(video)
		if err == nil {
			err = cache.Set(task.Id, string(marshal))
		}
	}

	if err != nil {
		log.Error("Failed set status processed" + err.Error())
	}
}
