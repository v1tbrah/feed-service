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

func TestCache_AddPostToUser(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		userID      int64
		inputPosts  []model.Post
		wantOutputs []model.Post
	}{
		{
			name:   "add 1 post",
			userID: 1,
			inputPosts: []model.Post{
				{
					ID:          1,
					UserID:      1,
					Description: "TestDescription1",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(10, 0).UTC(),
				},
			},
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
			name:   "add 2 posts",
			userID: 1,
			inputPosts: []model.Post{
				{
					ID:          1,
					UserID:      1,
					Description: "TestDescription1",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(10, 0).UTC(),
				},
				{
					ID:          2,
					UserID:      1,
					Description: "TestDescription2",
					HashtagsID:  []int64{3, 4},
					CreatedAt:   time.Unix(20, 0).UTC(),
				},
			},
			wantOutputs: []model.Post{
				{
					ID:          2,
					UserID:      1,
					Description: "TestDescription2",
					HashtagsID:  []int64{3, 4},
					CreatedAt:   time.Unix(20, 0).UTC(),
				},
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
			name:   "add 11 posts (1 more then maxPostToUser)",
			userID: 1,
			inputPosts: []model.Post{
				{
					ID:          1,
					UserID:      1,
					Description: "TestDescription1",
					HashtagsID:  []int64{1, 2},
					CreatedAt:   time.Unix(10, 0).UTC(),
				},
				{
					ID:          2,
					UserID:      1,
					Description: "TestDescription2",
					HashtagsID:  []int64{3, 4},
					CreatedAt:   time.Unix(20, 0).UTC(),
				},
				{
					ID:          3,
					UserID:      1,
					Description: "TestDescription3",
					HashtagsID:  []int64{5, 6},
					CreatedAt:   time.Unix(30, 0).UTC(),
				},
				{
					ID:          4,
					UserID:      1,
					Description: "TestDescription4",
					HashtagsID:  []int64{7, 8},
					CreatedAt:   time.Unix(40, 0).UTC(),
				},
				{
					ID:          5,
					UserID:      1,
					Description: "TestDescription5",
					HashtagsID:  []int64{9, 10},
					CreatedAt:   time.Unix(50, 0).UTC(),
				},
				{
					ID:          6,
					UserID:      1,
					Description: "TestDescription6",
					HashtagsID:  []int64{11, 12},
					CreatedAt:   time.Unix(60, 0).UTC(),
				},
				{
					ID:          7,
					UserID:      1,
					Description: "TestDescription7",
					HashtagsID:  []int64{13, 14},
					CreatedAt:   time.Unix(70, 0).UTC(),
				},
				{
					ID:          8,
					UserID:      1,
					Description: "TestDescription8",
					HashtagsID:  []int64{15, 16},
					CreatedAt:   time.Unix(80, 0).UTC(),
				},
				{
					ID:          9,
					UserID:      1,
					Description: "TestDescription9",
					HashtagsID:  []int64{17, 18},
					CreatedAt:   time.Unix(90, 0).UTC(),
				},
				{
					ID:          10,
					UserID:      1,
					Description: "TestDescription10",
					HashtagsID:  []int64{19, 20},
					CreatedAt:   time.Unix(100, 0).UTC(),
				},
				{
					ID:          11,
					UserID:      1,
					Description: "TestDescription11",
					HashtagsID:  []int64{21, 22},
					CreatedAt:   time.Unix(110, 0).UTC(),
				},
			},
			wantOutputs: []model.Post{
				{
					ID:          11,
					UserID:      1,
					Description: "TestDescription11",
					HashtagsID:  []int64{21, 22},
					CreatedAt:   time.Unix(110, 0).UTC(),
				},
				{
					ID:          10,
					UserID:      1,
					Description: "TestDescription10",
					HashtagsID:  []int64{19, 20},
					CreatedAt:   time.Unix(100, 0).UTC(),
				},
				{
					ID:          9,
					UserID:      1,
					Description: "TestDescription9",
					HashtagsID:  []int64{17, 18},
					CreatedAt:   time.Unix(90, 0).UTC(),
				},
				{
					ID:          8,
					UserID:      1,
					Description: "TestDescription8",
					HashtagsID:  []int64{15, 16},
					CreatedAt:   time.Unix(80, 0).UTC(),
				},
				{
					ID:          7,
					UserID:      1,
					Description: "TestDescription7",
					HashtagsID:  []int64{13, 14},
					CreatedAt:   time.Unix(70, 0).UTC(),
				},
				{
					ID:          6,
					UserID:      1,
					Description: "TestDescription6",
					HashtagsID:  []int64{11, 12},
					CreatedAt:   time.Unix(60, 0).UTC(),
				},
				{
					ID:          5,
					UserID:      1,
					Description: "TestDescription5",
					HashtagsID:  []int64{9, 10},
					CreatedAt:   time.Unix(50, 0).UTC(),
				},
				{
					ID:          4,
					UserID:      1,
					Description: "TestDescription4",
					HashtagsID:  []int64{7, 8},
					CreatedAt:   time.Unix(40, 0).UTC(),
				},
				{
					ID:          3,
					UserID:      1,
					Description: "TestDescription3",
					HashtagsID:  []int64{5, 6},
					CreatedAt:   time.Unix(30, 0).UTC(),
				},
				{
					ID:          2,
					UserID:      1,
					Description: "TestDescription2",
					HashtagsID:  []int64{3, 4},
					CreatedAt:   time.Unix(20, 0).UTC(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tHelperInitEmptyCache(t)
			for _, post := range tt.inputPosts {
				if post.UserID != tt.userID {
					t.Fatal("invalid test case: post.UserID != tt.userID")
				}

				post.UserID = tt.userID
				if err := c.AddPostToUser(ctx, tt.userID, post); err != nil {
					t.Fatalf("AddPostToUser() error = %v", err)
				}
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
			require.Equal(t, tt.wantOutputs, gotOutputs)
		})
	}
}
