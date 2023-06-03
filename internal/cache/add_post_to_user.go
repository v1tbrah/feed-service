package cache

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"gitlab.com/pet-pr-social-network/feed-service/internal/model"
)

const maxPostToUser = 10

func (c *Cache) AddPostToUser(ctx context.Context, userID int64, post model.Post) error {
	postData, err := json.Marshal(post)
	if err != nil {
		return errors.Wrapf(err, "json.Marshal, post: %+v", post)
	}

	userIDKey := strconv.Itoa(int(userID))

	if _, err = c.cli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		var errPipe error
		defer func() {
			if errPipe != nil {
				pipe.Discard()
			}
		}()

		if errPipe = pipe.LPush(ctx, userIDKey, postData).Err(); errPipe != nil {
			return errors.Wrapf(errPipe, "lpush, userID %s", userIDKey)
		}

		if errPipe = pipe.LTrim(ctx, userIDKey, 0, maxPostToUser-1).Err(); errPipe != nil {
			return errors.Wrapf(errPipe, "ltrim, userID %s", userIDKey)
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "pipe")
	}

	return nil
}

func (c *Cache) AddPostsToUser(ctx context.Context, userID int64, posts []model.Post) error {
	userIDKey := strconv.Itoa(int(userID))

	if _, err := c.cli.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		var errPipe error
		defer func() {
			if errPipe != nil {
				pipe.Discard()
			}
		}()

		for _, post := range posts {
			postData, err := json.Marshal(post)
			if err != nil {
				return errors.Wrapf(err, "json.Marshal, post: %+v", post)
			}

			if errPipe = pipe.LPush(ctx, userIDKey, postData).Err(); errPipe != nil {
				return errors.Wrapf(errPipe, "lpush, userID %s", userIDKey)
			}
		}

		if errPipe = pipe.LTrim(ctx, userIDKey, 0, maxPostToUser-1).Err(); errPipe != nil {
			return errors.Wrapf(errPipe, "ltrim, userID %s", userIDKey)
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "pipe")
	}

	return nil
}
