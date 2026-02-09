package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"
)

// VoteForPost 为帖子投票logic
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(*p.Direction))
}
