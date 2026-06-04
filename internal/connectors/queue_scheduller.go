package connectors

import (
	"time"

	queue_scheduler_pb "github.com/SV1Stail/tg-project-protos/gen/go/queue_scheduler/queue_scheduler"
	"google.golang.org/grpc"
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
	conn, err := grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &QueueschedulerClient{
		conn:                 conn,
		QueueschedulerClient: queue_scheduler_pb.NewQueueschedulerClient(conn),
	}, nil
}

func (c *QueueschedulerClient) Close() error {
	return c.conn.Close()
}
