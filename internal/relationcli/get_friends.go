package relationcli

import (
	"context"

	"github.com/pkg/errors"

	"gitlab.com/pet-pr-social-network/relation-service/rpbapi"
)

func (r *RelationCli) GetFriends(ctx context.Context, userID int64) ([]int64, error) {
	resp, err := r.cli.GetFriends(ctx, &rpbapi.GetFriendsRequest{UserID: userID})
	if err != nil {
		return nil, errors.Wrapf(err, "cli.GetFriends, userID: %d", userID)
	}

	if resp == nil {
		return nil, errors.Errorf("nil resp from cli.GetFriends, userID: %d", userID)
	}

	result := make([]int64, 0, len(resp.GetFriends()))
	for _, f := range resp.GetFriends() {
		result = append(result, f)
	}

	return result, nil
}
