package eventreader

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type MsgFriendRemoved struct {
	UserID   int64
	FriendID int64
}

func (mr *Reader) startProcessEventFriendRemoved() {
	if mr == nil {
		return
	}

	for {
		mv, err := mr.readerFriendRemoved.ReadMessage(mr.liveCtx)

		if err != nil {
			log.Warn().Err(err).Msg("read message friend removed")
			continue
		}

		var msg MsgFriendRemoved
		if err = json.Unmarshal(mv.Value, &msg); err != nil {
			log.Error().Err(err).Msg("json.Unmarshal friend removed message value")
			continue
		}

		if err = mr.cache.RemovePostsByUserID(mr.liveCtx, msg.UserID, msg.FriendID); err != nil {
			log.Error().Err(err).Int64("userID", msg.UserID).Int64("userIDWithWhichPostsNeedRemove", msg.FriendID).Msg("cache.RemovePostsByUserID")
			continue
		}

		if err = mr.cache.RemovePostsByUserID(mr.liveCtx, msg.FriendID, msg.UserID); err != nil {
			log.Error().Err(err).Int64("userID", msg.FriendID).Int64("userIDWithWhichPostsNeedRemove", msg.UserID).Msg("cache.RemovePostsByUserID")
			continue
		}
	}
}
