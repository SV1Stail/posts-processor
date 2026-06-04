package app

import (
	"sync"

	"github.com/SV1Stail/posts-processor/internal/connectors"
	"github.com/SV1Stail/posts-processor/internal/db"
)

type PostProcessor struct {
	QueueSchedulerClient *connectors.QueueschedulerClient
	DB                   *db.DB
	stopCh               chan struct{}
	mu                   *sync.Mutex
}

func NewQueueSchedulerService(queueSchedulerClient *connectors.QueueschedulerClient, db *db.DB) *PostProcessor {
	return &PostProcessor{
		stopCh:               make(chan struct{}),
		mu:                   &sync.Mutex{},
		QueueSchedulerClient: queueSchedulerClient,
		DB:                   db,
	}
}

func (qs *PostProcessor) Close() {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	qs.QueueSchedulerClient.Close()
	qs.DB.Close()
	close(qs.stopCh)
}
