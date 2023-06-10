package relationcli

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/v1tbrah/feed-service/config"

	"github.com/v1tbrah/relation-service/rpbapi"
)

type RelationCli struct {
	conn *grpc.ClientConn
	cli  rpbapi.RelationServiceClient
}

func New(cfg config.RelationCli) (*RelationCli, error) {
	conn, err := grpc.Dial(net.JoinHostPort(cfg.Host, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial")
	}

	return &RelationCli{
		conn: conn,
		cli:  rpbapi.NewRelationServiceClient(conn),
	}, nil
}

func (r *RelationCli) Close(ctx context.Context) (err error) {
	closed := make(chan struct{})

	go func() {
		if err = r.conn.Close(); err != nil {
			err = errors.Wrap(err, "conn.Close")
			closed <- struct{}{}
			return
		}

		closed <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-closed:
		return err
	}
}
