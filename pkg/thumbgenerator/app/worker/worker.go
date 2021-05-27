package worker

import (
	"github.com/bearname/videohost/pkg/common/caching"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/thumbgenerator/app/provider"
	"github.com/bearname/videohost/pkg/thumbgenerator/app/subscriber"
	"github.com/bearname/videohost/pkg/thumbgenerator/domain/model"
	log "github.com/sirupsen/logrus"
	"sync"
)

const WorkersCount = 3

func PoolOfWorker(stopChan chan struct{}, db database.Connector, cache caching.Cache) *sync.WaitGroup {
	var waitGroup sync.WaitGroup
	tasksChan := provider.RunTaskProvider(stopChan, db)
	for i := 0; i < WorkersCount; i++ {
		go func(i int) {
			waitGroup.Add(1)
			worker(tasksChan, db, cache, i)
			waitGroup.Done()
		}(i)
	}
	return &waitGroup
}

func worker(tasksChan <-chan *model.Task, db database.Connector, cache caching.Cache, name int) {
	log.Printf("start worker %v\n", name)
	for task := range tasksChan {
		log.Printf("start processing video with id %v on worker %v\n", task.Id, name)
		subscriber.HandleTask(task, db, cache)
		log.Printf("end processing video with id %v on worker %v\n", task.Id, name)
	}
	log.Printf("stop worker %v\n", name)
}
