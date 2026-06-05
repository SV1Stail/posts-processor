package main

import (
	"context"

	"github.com/SV1Stail/posts-processor/internal/app"
	"github.com/SV1Stail/posts-processor/internal/connectors"
	"github.com/SV1Stail/posts-processor/internal/db"
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
	mistralClient := connectors.NewMistralClient(&connectors.NewMistralClientRequest{})

	app := app.NewQueueSchedulerService(queueschedulerClient, db, mistralClient)
	defer app.Close()

	app.ProcessPosts(ctx)
}
