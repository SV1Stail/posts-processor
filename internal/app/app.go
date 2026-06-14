package app

import (
	"context"
	"os"
	"sync"

	"github.com/SV1Stail/posts-processor/internal/connectors"
	"github.com/SV1Stail/posts-processor/internal/db"
	post_processor_pb "github.com/SV1Stail/tg-project-protos/gen/go/posts_processor"
)

type PostProcessor struct {
	post_processor_pb.UnimplementedPostsProcessorServer
	QueueSchedulerClient *connectors.QueueschedulerClient
	DB                   *db.DB
	LlmClient            LLM
	Port                 string
	stopCh               chan struct{}
	mu                   *sync.Mutex
}

type LLM interface {
	// without chat_id
	NewChat(ctx context.Context, systemPrompt, userMessage string) (string, error)
}

func NewPostProcessorService(queueSchedulerClient *connectors.QueueschedulerClient, db *db.DB, llm LLM) *PostProcessor {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8090"
	}

	return &PostProcessor{
		QueueSchedulerClient: queueSchedulerClient,
		DB:                   db,
		LlmClient:            llm,
		Port:                 port,
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
