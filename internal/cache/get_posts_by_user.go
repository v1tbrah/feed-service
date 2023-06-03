package cache

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"

	"gitlab.com/pet-pr-social-network/feed-service/internal/model"
)

func (c *Cache) GetPostsByUserID(ctx context.Context, userID int64) ([]model.Post, error) {
	res, err := c.cli.LRange(ctx, strconv.Itoa(int(userID)), 0, maxPostToUser-1).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "lrange, userID %d", userID)
	}

	posts := make([]model.Post, 0, len(res))
	for _, postStr := range res {
		var post model.Post
		if err = json.Unmarshal([]byte(postStr), &post); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal, postStr: %s", postStr)
		}
		posts = append(posts, post)
	}

	return posts, nil
}
