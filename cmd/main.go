package main

import (
	"context"

	"github.com/SV1Stail/posts-processor/internal/app"
	"github.com/SV1Stail/posts-processor/internal/connectors"
	"github.com/SV1Stail/posts-processor/internal/db"
	"github.com/gage-technologies/mistral-go"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := db.MustNewDB(ctx)
	queueschedulerClient, err := connectors.NewQueueSchedulerClient(&connectors.QueueschedulerConfig{})
	if err != nil {
		log.Err(err).Ctx(ctx).Msg("Failed to create queue-scheduler client")
	}

	mistralClient := mistral.NewMistralClientDefault("your-api-key")

	app := app.NewQueueSchedulerService(queueschedulerClient, db)
	defer app.Close()

	// Example: Using Chat Completions
	chatRes, err := mistralClient.Chat("mistral-tiny", []mistral.ChatMessage{{Content: "Hello, world!", Role: mistral.RoleUser}}, nil)
	if err != nil {
		log.Err(err).Msg("Error getting chat completion")
	}
	log.Printf("Chat completion: %+v\n", chatRes)

}
