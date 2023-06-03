package eventreader

import (
	"context"

	"gitlab.com/pet-pr-social-network/feed-service/internal/model"
)

//go:generate mockery --name Cache
type Cache interface {
	AddPostToUser(ctx context.Context, userID int64, post model.Post) error
	AddPostsToUser(ctx context.Context, userID int64, posts []model.Post) error
	RemovePostFromUser(ctx context.Context, userID, postID int64) error
	RemovePostsByUserID(ctx context.Context, userID int64, userIDWithWhichPostsNeedRemove int64) error
}

//go:generate mockery --name RelationCli
type RelationCli interface {
	GetFriends(ctx context.Context, userID int64) ([]int64, error)
}

//go:generate mockery --name PostCli
type PostCli interface {
	GetPostsByUserID(ctx context.Context, userID int64) ([]model.Post, error)
}
