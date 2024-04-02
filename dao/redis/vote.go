package redis

import (
	"context"
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

// CreatePost redis存储帖子信息 使用hash存储帖子信息
func CreatePost(postID int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 创建帖子时间和帖子分数要满足事务性
	pipeline := rdb.TxPipeline()

	// 帖子时间
	now := float64(time.Now().Unix()) // Unix时间戳: 从 Unix纪元（1970年1月1日 UTC）到当前时间的总秒数
	_, err = pipeline.ZAdd(ctx, KeyPostTimeZSet, redis.Z{
		Score:  now,
		Member: postID,
	}).Result()

	// 帖子分数
	_, err = pipeline.ZAdd(ctx, KeyPostScoreZSet, redis.Z{
		Score:  now,
		Member: postID,
	}).Result()

	// 执行pipeline
	_, err = pipeline.Exec(ctx)
	return err
}

// VoteForPost	为帖子投票
func VoteForPost(userID, postID string, value float64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 1. 判断投票限制
	// 从Redis获取帖子发布时间
	postTime := rdb.ZScore(ctx, KeyPostTimeZSet, postID).Val()
	// 时间作差，判断是否超过一周
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrorVoteTimeExpire
	}

	// 2. 更新帖子分数
	// 先查询当前用户给当前帖子的投票记录
	oldValue := rdb.ZScore(ctx, KeyPostVotedZSetPrefix+postID, userID).Val()
	// 更新：如果这一次投票的值和之前保存的值一致，提示不允许重复投票
	if value == oldValue {
		return ErrorVoteRepeated
	}

	// 2和3需要放到一个pipeline事务操作
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, KeyPostScoreZSet, (value-oldValue)*scorePerVote, postID)

	// 3. 记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(ctx, KeyPostVotedZSetPrefix+postID, userID)
	} else {
		pipeline.ZAdd(ctx, KeyPostVotedZSetPrefix+postID, redis.Z{ // 更新或插入新的ZSet
			Score:  value,
			Member: userID,
		})
	}
	_, err = pipeline.Exec(ctx)
	return err
}
