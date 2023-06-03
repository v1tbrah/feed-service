package eventreader

import (
	"encoding/json"

	"github.com/rs/zerolog/log"

	"gitlab.com/pet-pr-social-network/feed-service/internal/model"
)

func (mr *Reader) startProcessEventPostCreated() {
	if mr == nil {
		return
	}

	for {
		m, err := mr.readerPostCreated.ReadMessage(mr.liveCtx)

		if err != nil {
			log.Warn().Err(err).Msg("read message post created")
			continue
		}

		var post model.Post
		if err = json.Unmarshal(m.Value, &post); err != nil {
			log.Error().Err(err).Msg("json.Unmarshal post created message value")
			continue
		}

		if err = mr.cache.AddPostToUser(mr.liveCtx, post.UserID, post); err != nil {
			log.Error().Err(err).Interface("post", post).Int64("userID", post.UserID).Msg("cache.AddPostToUser")
			continue
		}

		friends, err := mr.relationCli.GetFriends(mr.liveCtx, post.UserID)
		if err != nil {
			log.Error().Err(err).Int64("userID", post.UserID).Msg("relationCli.GetFriends")
			continue
		}

		for _, friendID := range friends {
			if err = mr.cache.AddPostToUser(mr.liveCtx, friendID, post); err != nil {
				log.Error().Err(err).Interface("post", post).Int64("userID", friendID).Msg("cache.AddPostToUser")
				continue
			}
		}
	}
}
