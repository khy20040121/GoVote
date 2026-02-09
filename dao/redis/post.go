package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// GetIDsFromKey 根据key获得ids
func GetIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1

	return client.ZRevRange(key, start, end).Result()
}

// CreatePost 初始化redis中的帖子
func CreatePost(p *models.Post) error {
	// 封转成一个事务来做
	pipe := client.TxPipeline()
	// 初始化分数
	pipe.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Member: p.ID,
		Score:  float64(time.Now().Unix()),
	})
	// 初始化时间
	pipe.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Member: p.ID,
		Score:  float64(time.Now().Unix()),
	})

	// 把帖子的社区id加入到redis中去
	pipe.ZAdd(getRedisKey(KeyCommunitySetPF+strconv.Itoa(int(p.CommunityID))), redis.Z{
		Member: p.ID,
		Score:  1,
	})
	_, err := pipe.Exec()
	return err
}

// GetPostIDsInOrder根据指定顺序获取帖子列表
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	key := getRedisKey(KeyPostScoreZSet)
	if p.Order == models.OrderTime {
		key = getRedisKey(KeyPostTimeZSet)
	}

	return GetIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteList获取帖子的赞成票数
func GetPostVoteList(ids []string) (data []int64, err error) {
	pipe := client.Pipeline()

	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipe.ZCount(key, "1", "1")
		pipe.ZCount(key, "-1", "-1")
	}
	cmders, err := pipe.Exec()
	if err != nil {
		return
	}

	data = make([]int64, 0, len(ids))
	for i := 0; i < len(cmders); i += 2 {
		v1 := cmders[i].(*redis.IntCmd).Val()
		v2 := cmders[i+1].(*redis.IntCmd).Val()
		data = append(data, v1-v2)
	}
	return
}

// GetCommunityPostIDsInOrder按社区获取帖子的ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	key := orderKey + ":" + strconv.Itoa(int(p.CommunityID))

	// 如果不存在需要计算
	pipe := client.Pipeline()
	pipe.ZInterStore(key, redis.ZStore{
		Aggregate: "MAX",
	}, cKey, orderKey)
	pipe.Expire(key, 60*time.Second)
	_, err := pipe.Exec()
	if err != nil {
		return nil, err
	}

	return GetIDsFromKey(key, p.Page, p.Size)
}
