package relationcli

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.com/pet-pr-social-network/feed-service/config"
	"gitlab.com/pet-pr-social-network/relation-service/rpbapi"
)

type RelationCli struct {
	conn *grpc.ClientConn
	cli  rpbapi.RelationServiceClient
}

func New(cfg config.RelationCli) (*RelationCli, error) {
	conn, err := grpc.Dial(net.JoinHostPort(cfg.ServHost, cfg.ServPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %w", err)
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
func (r *RelationCli) GetFriends(ctx context.Context, userID int64) ([]int64, error) {
	resp, err := r.cli.GetFriends(ctx, &rpbapi.GetFriendsRequest{UserID: userID})
	if err != nil {
		return nil, errors.Wrapf(err, "cli.GetFriends, userID: %d", userID)
	}

	if resp == nil {
		return nil, errors.Errorf("nil resp from cli.GetFriends, userID: %d", userID)
	}

	result := make([]int64, 0, len(resp.GetFriends()))
	for _, f := range resp.GetFriends() {
		result = append(result, f)
	}

	return result, nil
}
