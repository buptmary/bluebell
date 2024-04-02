package redis

import (
	"bluebell/models"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// getIDsFromKey 按照分数从大到小的顺序查询指定数量的元素
func getIDsFromKey(key string, page, size int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	start := (page - 1) * size
	end := start + size - 1
	// ZRevRange 按照分数从大到小的顺序查询指定数量的元素
	return rdb.ZRevRange(ctx, key, start, end).Result()
}

// GetPostIDsInOrder 升级版投票列表接口
// 按创建时间排序 或者 按照 分数排序 (查询出的ids已经根据order从大到小排序)
func GetPostIDsInOrder(p *models.PostListForm) ([]string, error) {
	// 从redis获取id
	// 1. 根据用户请求中携带的order参数确定要查询的redis key
	key := KeyPostTimeZSet            // 默认是时间
	if p.Order == models.OrderScore { // 按照分数请求
		key = KeyPostScoreZSet
	}
	// 2. 确定查询的索引起始点
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	//for _, id := range ids {
	//	key := KeyPostVotedZSetPrefix + id
	//	// 查找key中分数为1的元素数量 -> 统计每篇帖子的赞成票的数量
	//	v := rdb.ZCount(ctx, key, "1", "1").Val()
	//	data = append(data, v)
	//}
	// 使用pipeline一次发送多条命令减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := KeyPostVotedZSetPrefix + id
		// 查找key中分数为1的元素数量 -> 统计每篇帖子的赞成票的数量
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return data, nil
}
