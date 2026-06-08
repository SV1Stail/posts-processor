package main

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/SV1Stail/posts-processor/internal/app"
	"github.com/SV1Stail/posts-processor/internal/connectors"
	"github.com/SV1Stail/posts-processor/internal/db"
	post_processor_pb "github.com/SV1Stail/tg-project-protos/gen/go/posts_processor"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := db.MustNewDB(ctx)
	queueSchedulerClient, err := connectors.NewQueueSchedulerClient(&connectors.QueueschedulerConfig{
		Address: "queue-scheduler:8090",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		log.Err(err).Ctx(ctx).Msg("Failed to create queue-scheduler client")
	}
	mistralClient := connectors.NewMistralClient(&connectors.NewMistralClientRequest{})

	app := app.NewPostProcessorService(queueSchedulerClient, db, mistralClient)
	defer app.Close()

	grpcServer := grpc.NewServer()
	post_processor_pb.RegisterPostsProcessorServer(grpcServer, app)

	lis, err := net.Listen("tcp", ":"+app.Port)
	if err != nil {
		log.Err(err).Ctx(ctx).Msg("failed to listen")
	}

	go app.Workers(ctx)

	if err := grpcServer.Serve(lis); err != nil {
		log.Err(err).Ctx(ctx).Msg("failed to serve")
		os.Exit(1)
	}

	grpcServer.GracefulStop()
}
