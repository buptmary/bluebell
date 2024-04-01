package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoteRepeated   = errors.New("不允许重复投票")
)

// VoteForPost	为帖子投票
func VoteForPost(userID, postID string, value float64) (err error) {
	// 1. 判断投票限制
	// 从Redis获取帖子发布时间
	postTime := rdb.ZScore(ctx, KeyPostTimeZSet, postID).Val()
	// 时间作差，判断是否超过一周
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrorVoteTimeExpire
	}

	// 2. 更新帖子分数
	// 先查询当前用户给当前帖子的投票记录
	old_value := rdb.ZScore(ctx, KeyPostVotedZSetPrefix+postID, userID).Val()
	// 更新：如果这一次投票的值和之前保存的值一致，提示不允许重复投票
	if value == old_value {
		return ErrorVoteRepeated
	}
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, KeyPostScoreZSet, (value-old_value)*scorePerVote, postID)

	// 3. 记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(ctx, KeyPostVotedZSetPrefix+postID, userID)
	} else {
		pipeline.ZAdd(ctx, KeyPostVotedZSetPrefix+postID, redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err = pipeline.Exec(ctx)
	return err
}
