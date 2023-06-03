package cache

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	
	"gitlab.com/pet-pr-social-network/feed-service/internal/model"
)

func (c *Cache) RemovePostFromUser(ctx context.Context, userID int64, postID int64) error {
	userIDKey := strconv.Itoa(int(userID))

	posts, err := c.cli.LRange(ctx, strconv.Itoa(int(userID)), 0, maxPostToUser-1).Result()
	if err != nil {
		return errors.Wrapf(err, "lrange, userID %d", userID)
	}

	for _, postStr := range posts {
		var post model.Post
		if err = json.Unmarshal([]byte(postStr), &post); err != nil {
			return errors.Wrapf(err, "json.Unmarshal, postStr: %s", postStr)
		}

		if post.ID == postID {
			if err = c.cli.LRem(ctx, userIDKey, 1, []byte(postStr)).Err(); err != nil {
				return errors.Wrapf(err, "cli.LRem, userID (%d), post (%+v)", userID, post)
			}
			return nil
		}
	}

	return ErrPostNotFound
}

func (c *Cache) RemovePostsByUserID(ctx context.Context, userID int64, userIDWithWhichPostsNeedRemove int64) error {
	userIDKey := strconv.Itoa(int(userID))

	posts, err := c.cli.LRange(ctx, strconv.Itoa(int(userID)), 0, maxPostToUser-1).Result()
	if err != nil {
		return errors.Wrapf(err, "lrange, userID %d", userID)
	}

	for _, postStr := range posts {
		var post model.Post
		if err = json.Unmarshal([]byte(postStr), &post); err != nil {
			return errors.Wrapf(err, "json.Unmarshal, postStr: %s", postStr)
		}

		if post.UserID == userIDWithWhichPostsNeedRemove {
			if err = c.cli.LRem(ctx, userIDKey, 1, []byte(postStr)).Err(); err != nil {
				return errors.Wrapf(err, "cli.LRem, userID (%d), post (%+v)", userID, post)
			}
			return nil
		}
	}

	return nil
}
