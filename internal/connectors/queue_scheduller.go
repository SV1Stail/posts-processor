package connectors

import (
	"context"
	"fmt"
	"time"

	queue_scheduler_pb "github.com/SV1Stail/tg-project-protos/gen/go/queue_scheduler/queue_scheduler"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type QueueschedulerClient struct {
	conn *grpc.ClientConn
	queue_scheduler_pb.QueueschedulerClient
}

type QueueschedulerConfig struct {
	Address string        // "publisher:9090" или "localhost:9090"
	Timeout time.Duration // например: 30 * time.Second
}

func NewQueueSchedulerClient(cfg *QueueschedulerConfig) (*QueueschedulerClient, error) {
	// if cfg.Address == "" {
	// 	cfg.Address = os.Getenv("QUEUE_SCHEDULER_ADDR")
	// 	if cfg.Address == "" {
	// 		cfg.Address = "queue-scheduler:8090"
	// 	}
	// }

	conn, err := grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := queue_scheduler_pb.NewQueueschedulerClient(conn)

	return &QueueschedulerClient{
		conn:                 conn,
		QueueschedulerClient: client,
	}, nil
}

func (c *QueueschedulerClient) Close() error {
	return c.conn.Close()
}

func (c *QueueschedulerClient) CheckConnectionState(ctx context.Context) error {
	if c == nil {
		log.Err(fmt.Errorf("client is nil")).Ctx(ctx).Msg("QueueSchedulerClient is nil")

		return fmt.Errorf("client is nil")
	}
	state := c.conn.GetState()

	if state == connectivity.Idle {
		log.Warn().Ctx(ctx).
			Str("state", state.String()).
			Msg("QueueSchedulerClient connection is IDLE")

		return nil
	}

	if state != connectivity.Ready {
		log.Err(fmt.Errorf("connection is not ready")).Ctx(ctx).
			Str("state", state.String()).
			Msg("QueueSchedulerClient connection is not ready")

		return fmt.Errorf("connection is not ready")
	}

	log.Info().Msg("CheckConnectionState SUCCESS")

	return nil

}
