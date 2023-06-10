package postcli

import (
	"context"

	"github.com/pkg/errors"

	"github.com/v1tbrah/feed-service/internal/model"

	"github.com/v1tbrah/post-service/ppbapi"
)

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
