package eventreader

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type MsgFriendAdded struct {
	UserID   int64
	FriendID int64
}

func (mr *Reader) startProcessEventFriendAdded() {
	if mr == nil {
		return
	}

	for {
		mv, err := mr.readerFriendAdded.ReadMessage(mr.liveCtx)

		if err != nil {
			log.Warn().Err(err).Msg("read message friend added")
			continue
		}

		var msg MsgFriendAdded
		if err = json.Unmarshal(mv.Value, &msg); err != nil {
			log.Error().Err(err).Msg("json.Unmarshal friend added message value")
			continue
		}

		postsForUser, err := mr.postCli.GetPostsByUserID(mr.liveCtx, msg.UserID)
		if err != nil {
			log.Error().Err(err).Int64("userID", msg.UserID).Msg("postCli.GetPostsByUserID")
			continue
		}

		if err = mr.cache.AddPostsToUser(mr.liveCtx, msg.FriendID, postsForUser); err != nil {
			log.Error().Err(err).Int64("userID", msg.FriendID).Interface("posts", postsForUser).Msg("cache.AddPostsToUser")
			continue
		}

		postsForFriend, err := mr.postCli.GetPostsByUserID(mr.liveCtx, msg.FriendID)
		if err != nil {
			log.Error().Err(err).Int64("userID", msg.FriendID).Msg("postCli.GetPostsByUserID")
			continue
		}

		if err = mr.cache.AddPostsToUser(mr.liveCtx, msg.UserID, postsForFriend); err != nil {
			log.Error().Err(err).Int64("userID", msg.UserID).Interface("posts", postsForFriend).Msg("cache.AddPostsToUser")
			continue
		}
	}
}
