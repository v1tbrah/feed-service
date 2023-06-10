package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/v1tbrah/feed-service/fpbapi"
	"github.com/v1tbrah/feed-service/internal/api/mocks"
	"github.com/v1tbrah/feed-service/internal/model"
)

func TestAPI_GetFeed(t *testing.T) {
	tests := []struct {
		name            string
		mockCache       func(t *testing.T) *mocks.Cache
		req             *fpbapi.GetFeedRequest
		expectedResp    *fpbapi.GetFeedResponse
		wantErr         bool
		expectedErr     error
		expectedErrCode codes.Code
	}{
		{
			name: "OK",
			req:  &fpbapi.GetFeedRequest{UserID: 1},
			mockCache: func(t *testing.T) *mocks.Cache {
				testCache := mocks.NewCache(t)
				postsFromCache := []model.Post{
					{
						ID:          1,
						UserID:      1,
						Description: "TestDescription1",
						HashtagsID:  []int64{1, 2},
						CreatedAt:   time.Unix(10, 20),
					},
					{
						ID:          2,
						UserID:      2,
						Description: "TestDescription2",
						HashtagsID:  []int64{3, 4},
						CreatedAt:   time.Unix(30, 40),
					},
				}
				testCache.On("GetPostsByUserID",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1)).
					Return(postsFromCache, nil).
					Once()
				return testCache
			},
			expectedResp: &fpbapi.GetFeedResponse{
				Posts: []*fpbapi.Post{
					{
						Id:          1,
						UserID:      1,
						Description: "TestDescription1",
						HashtagsID:  []int64{1, 2},
						CreatedAt:   timestamppb.New(time.Unix(10, 20)),
					},
					{
						Id:          2,
						UserID:      2,
						Description: "TestDescription2",
						HashtagsID:  []int64{3, 4},
						CreatedAt:   timestamppb.New(time.Unix(30, 40)),
					},
				},
			},
		},
		{
			name: "unexpected err on storage.GetPostsByUserID",
			req:  &fpbapi.GetFeedRequest{UserID: 1},
			mockCache: func(t *testing.T) *mocks.Cache {
				testCache := mocks.NewCache(t)
				testCache.On("GetPostsByUserID",
					mock.MatchedBy(func(ctx context.Context) bool { return true }), int64(1)).
					Return([]model.Post{}, errors.New("unexpected err")).
					Once()
				return testCache
			},
			wantErr:         true,
			expectedErr:     errors.New("unexpected"),
			expectedErrCode: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{cache: tt.mockCache(t)}
			resp, err := a.GetFeed(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrCode, status.Code(err))
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			}

			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
