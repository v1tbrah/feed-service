package eventreader

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type MsgPostDeleted struct {
	ID     int64
	UserID int64
}

func (mr *Reader) startProcessEventPostDeleted() {
	if mr == nil {
		return
	}

	for {
		mv, err := mr.readerPostDeleted.ReadMessage(mr.liveCtx)

		if err != nil {
			log.Warn().Err(err).Msg("read message post deleted")
			continue
		}

		var msg MsgPostDeleted
		if err = json.Unmarshal(mv.Value, &msg); err != nil {
			log.Error().Err(err).Msg("json.Unmarshal post deleted message value")
			continue
		}

		if err = mr.cache.RemovePostFromUser(mr.liveCtx, msg.UserID, msg.ID); err != nil {
			log.Error().Err(err).Int64("id", msg.ID).Int64("userID", msg.UserID).Msg("cache.RemovePostFromUser")
			continue
		}

		friends, err := mr.relationCli.GetFriends(mr.liveCtx, msg.UserID)
		if err != nil {
			log.Error().Err(err).Int64("userID", msg.UserID).Msg("relationCli.GetFriends")
			continue
		}

		for _, friendID := range friends {
			if err = mr.cache.RemovePostFromUser(mr.liveCtx, msg.ID, friendID); err != nil {
				log.Error().Err(err).Int64("id", msg.ID).Int64("userID", friendID).Msg("cache.RemovePostFromUser")
				continue
			}
		}
	}
}
