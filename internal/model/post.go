package model

import "time"

type Post struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Description string    `json:"description"`
	HashtagsID  []int64   `json:"hashtags_id"`
	CreatedAt   time.Time `json:"created_at"`
}
