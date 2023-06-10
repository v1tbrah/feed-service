package postcli

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/v1tbrah/feed-service/config"

	"github.com/v1tbrah/post-service/ppbapi"
)

type PostCli struct {
	conn *grpc.ClientConn
	cli  ppbapi.PostServiceClient
}

func New(cfg config.PostCli) (*PostCli, error) {
	conn, err := grpc.Dial(net.JoinHostPort(cfg.Host, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial")
	}

	return &PostCli{
		conn: conn,
		cli:  ppbapi.NewPostServiceClient(conn),
	}, nil
}

func (pcli *PostCli) Close(ctx context.Context) (err error) {
	closed := make(chan struct{})

	go func() {
		if err = pcli.conn.Close(); err != nil {
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
