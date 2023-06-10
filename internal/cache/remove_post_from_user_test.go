//go:build with_db

package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/v1tbrah/feed-service/internal/model"
)

func TestCache_RemovePostFromUser(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string

		userID int64
		postID int64

		input []model.Post

		wantResult []model.Post
		wantErr    bool
	}{
		{
			name:   "simple test",
			userID: 1,
			postID: 1,

			input: []model.Post{
				{
					ID:          1,
					UserID:      1,
					Description: "TestDescription1",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(1000, 0).UTC(),
				},
				{
					ID:          2,
					UserID:      1,
					Description: "TestDescription2",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(2000, 0).UTC(),
				},
			},

			wantResult: []model.Post{
				{
					ID:          2,
					UserID:      1,
					Description: "TestDescription2",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(2000, 0).UTC(),
				},
			},
		},
		{
			name:   "remove non-existent post",
			userID: 1,
			postID: 3,

			input: []model.Post{
				{
					ID:          1,
					UserID:      1,
					Description: "TestDescription1",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(1000, 0).UTC(),
				},
				{
					ID:          2,
					UserID:      1,
					Description: "TestDescription2",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(2000, 0).UTC(),
				},
			},

			wantResult: []model.Post{
				{
					ID:          2,
					UserID:      1,
					Description: "TestDescription2",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(2000, 0).UTC(),
				},
				{
					ID:          1,
					UserID:      1,
					Description: "TestDescription1",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(1000, 0).UTC(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tHelperInitEmptyCache(t)

			for _, post := range tt.input {
				postData, err := json.Marshal(post)
				if err != nil {
					t.Fatalf("Marshal() error = %v", err)
				}

				if err = c.cli.LPush(ctx, strconv.Itoa(int(post.UserID)), postData).Err(); err != nil {
					t.Fatalf("LPush() error = %v", err)
				}
			}

			err := c.RemovePostFromUser(ctx, tt.userID, tt.postID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			res, err := c.cli.LRange(ctx, strconv.Itoa(int(tt.userID)), 0, maxPostToUser-1).Result()
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}

			gotOutputs := make([]model.Post, 0, len(res))
			for _, post := range res {
				var tempGotPost model.Post
				if err = json.Unmarshal([]byte(post), &tempGotPost); err != nil {
					t.Fatalf("Unmarshal() error = %v", err)
				}
				gotOutputs = append(gotOutputs, tempGotPost)
			}
			assert.Equal(t, tt.wantResult, gotOutputs)
		})
	}
}
