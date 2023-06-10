package api

import (
	"context"

	"github.com/v1tbrah/feed-service/internal/model"
)

//go:generate mockery --name Cache
type Cache interface {
	GetPostsByUserID(ctx context.Context, userID int64) ([]model.Post, error)
}
