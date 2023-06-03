package api

import (
	"context"

	"gitlab.com/pet-pr-social-network/feed-service/internal/model"
)

//go:generate mockery --name Cache
type Cache interface {
	GetPostsByUserID(ctx context.Context, userID int64) ([]model.Post, error)
}
