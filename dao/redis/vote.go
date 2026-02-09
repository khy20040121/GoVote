package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

// VoteForPost 为帖子投票dao
func VoteForPost(userID, postID string, dir float64) error {
	// 1 判断帖子投票限制,帖子一周之内才能投票

	PostTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-PostTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 投票分数的计算
	odir := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	// 投票重复,返回错误
	if odir == dir {
		return ErrVoteRepeated
	}
	var op float64
	if dir > odir {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(dir - odir)

	// 2 记录投票数据
	// 3 更新帖子分数
	// 2和 3需要放到一个事物当中去执行
	pipe := client.TxPipeline()
	// 更新分数
	pipe.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	zap.L().Info("", zap.Float64("op", op), zap.Float64("diff", diff), zap.Float64("odir", odir),
		zap.Float64("score", op*diff*scorePerVote),
	)

	//更新投票情况
	if dir == 0 {
		pipe.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipe.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Member: userID,
			Score:  dir,
		})
	}
	_, err := pipe.Exec()
	return err
}

// GetPostVoteForUser 获取用户对帖子的投票记录
func GetPostVoteForUser(userID, postID string) (float64, error) {
	return client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Result()
}
