package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/v1tbrah/feed-service/fpbapi"
)

func (a *API) GetFeed(ctx context.Context, req *fpbapi.GetFeedRequest) (*fpbapi.GetFeedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, fpbapi.ErrEmptyRequest.Error())
	}

	posts, err := a.cache.GetPostsByUserID(ctx, req.GetUserID())
	if err != nil {
		log.Err(err).Msg("storage.GetPostsByUserID")
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := make([]*fpbapi.Post, 0, len(posts))
	for _, post := range posts {
		res = append(res, &fpbapi.Post{
			Id:          post.ID,
			UserID:      post.UserID,
			Description: post.Description,
			HashtagsID:  post.HashtagsID,
			CreatedAt:   timestamppb.New(post.CreatedAt),
		})
	}

	return &fpbapi.GetFeedResponse{Posts: res}, nil
}
