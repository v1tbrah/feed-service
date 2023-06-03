package postcli

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.com/pet-pr-social-network/feed-service/config"
	"gitlab.com/pet-pr-social-network/feed-service/internal/model"
	"gitlab.com/pet-pr-social-network/post-service/ppbapi"
)

type PostCli struct {
	conn *grpc.ClientConn
	cli  ppbapi.PostServiceClient
}

func New(cfg config.PostCli) (*PostCli, error) {
	conn, err := grpc.Dial(net.JoinHostPort(cfg.ServHost, cfg.ServPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %w", err)
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

func (pcli *PostCli) GetPostsByUserID(ctx context.Context, userID int64) ([]model.Post, error) {
	resp, err := pcli.cli.GetPostsByUserID(ctx, &ppbapi.GetPostsByUserIDRequest{UserID: userID})
	if err != nil {
		return nil, errors.Wrapf(err, "cli.GetPostsByUserID, userID: %d", userID)
	}

	if resp == nil {
		return nil, errors.Errorf("nil resp from cli.GetPostsByUserID, userID: %d", userID)
	}

	result := make([]model.Post, 0, len(resp.GetPosts()))
	for _, p := range resp.GetPosts() {
		result = append(result, model.Post{
			ID:          p.GetId(),
			UserID:      p.GetUserID(),
			Description: p.GetDescription(),
			HashtagsID:  p.GetHashtagsID(),
			CreatedAt:   p.GetCreatedAt().AsTime(),
		})
	}

	return result, nil
}
