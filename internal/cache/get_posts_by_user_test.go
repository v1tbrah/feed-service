//go:build with_db

package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/v1tbrah/feed-service/internal/model"
)

func TestCache_GetPostsByUserID(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		userID     int64
		inputPosts []model.Post

		getUserIDList []int64
		wantOutputs   []model.Post
	}{
		{
			name: "add 1 post | get 1 post",
			inputPosts: []model.Post{
				{
					ID:          1,
					UserID:      1,
					Description: "TestDescription1",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(10, 0).UTC(),
				},
			},

			getUserIDList: []int64{1},
			wantOutputs: []model.Post{
				{
					ID:          1,
					UserID:      1,
					Description: "TestDescription1",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(10, 0).UTC(),
				},
			},
		},
		{
			name: "add 1 post user1 | get 1 post user2",
			inputPosts: []model.Post{
				{
					ID:          1,
					UserID:      1,
					Description: "TestDescription1",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(10, 0).UTC(),
				},
			},

			getUserIDList: []int64{2},
			wantOutputs:   []model.Post{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tHelperInitEmptyCache(t)
			for _, post := range tt.inputPosts {
				postData, err := json.Marshal(post)
				if err != nil {
					t.Fatalf("Marshal() error = %v", err)
				}

				if err = c.cli.LPush(ctx, strconv.Itoa(int(post.UserID)), postData).Err(); err != nil {
					t.Fatalf("LPush() error = %v", err)
				}
			}

			gotOutputs := make([]model.Post, 0)

			for _, userID := range tt.getUserIDList {
				res, err := c.GetPostsByUserID(ctx, userID)
				if err != nil {
					t.Fatalf("Get() error = %v", err)
				}
				gotOutputs = append(gotOutputs, res...)
			}

			require.Equal(t, tt.wantOutputs, gotOutputs)
		})
	}
}
