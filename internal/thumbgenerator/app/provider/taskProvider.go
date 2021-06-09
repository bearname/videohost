package provider

import (
	"github.com/bearname/videohost/internal/common/db"
	"github.com/bearname/videohost/internal/thumbgenerator/app/publisher"
	"github.com/bearname/videohost/internal/thumbgenerator/domain/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func RunTaskProvider(stopChan chan struct{}, db db.Connector) <-chan *model.Task {
	resultChan := make(chan *model.Task)
	stopTaskProviderChan := make(chan struct{})
	taskProviderChan := runTaskPublisher(stopTaskProviderChan, db)
	onStop := func() {
		stopTaskProviderChan <- struct{}{}
		close(resultChan)
	}
	go func() {
		for {
			select {
			case <-stopChan:
				onStop()
				return
			case task := <-taskProviderChan:
				select {
				case <-stopChan:
					onStop()
					return
				case resultChan <- task:
				}
			}
		}
	}()
	return resultChan
}

func runTaskPublisher(stopChan chan struct{}, db db.Connector) <-chan *model.Task {
	tasksChan := make(chan *model.Task)
	go func() {
		for {
			select {
			case <-stopChan:
				close(tasksChan)
				return
			default:
			}
			if task := publisher.PublishTask(db); task != nil {
				log.Printf("got the task %v\n", task)
				tasksChan <- task
			} else {
				log.Info("no task for processing, start waiting")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return tasksChan
}
