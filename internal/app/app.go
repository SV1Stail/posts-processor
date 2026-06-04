package app

import (
	"context"
	"sync"

	"github.com/SV1Stail/posts-processor/internal/connectors"
	"github.com/SV1Stail/posts-processor/internal/db"
)

type PostProcessor struct {
	QueueSchedulerClient *connectors.QueueschedulerClient
	DB                   *db.DB
	LlmClient            LLM
	stopCh               chan struct{}
	mu                   *sync.Mutex
}

type LLM interface {
	// without chat_id
	NewChat(ctx context.Context, systemPrompt, userMessage string) (string, error)
}

func NewQueueSchedulerService(queueSchedulerClient *connectors.QueueschedulerClient, db *db.DB, llm LLM) *PostProcessor {
	return &PostProcessor{
		QueueSchedulerClient: queueSchedulerClient,
		DB:                   db,
		LlmClient:            llm,
		stopCh:               make(chan struct{}),
		mu:                   &sync.Mutex{},
	}
}

func (qs *PostProcessor) Close() {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	qs.QueueSchedulerClient.Close()
	qs.DB.Close()
	close(qs.stopCh)
}
