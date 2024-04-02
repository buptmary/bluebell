package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"fmt"
	"go.uber.org/zap"
)

// 投票功能

/*
投票算法：http://www.ruanyifeng.com/blog/2012/03/ranking_algorithm_reddit.html
本项目使用简化版的投票分数
投一票就加432分 86400/200 -> 200张赞成票就可以给帖子在首页续1天  -> 《redis实战》
*/

/* PostVote 为帖子投票
投票分为四种情况：1.投赞成票(1) 2.投反对票(-1) 3.取消投票(0) 4.反转投票

记录文章参与投票的人
更新文章分数：赞成票要加分；反对票减分

v=1时，有两种情况
	1.之前没投过票，现在要投赞成票		--> 更新分数和投票记录		差值的绝对值：1  +432
	2.之前投过反对票，现在要改为赞成票	--> 更新分数和投票记录		差值的绝对值：2  +432*2
v=0时，有两种情况
	1.之前投过反对票，现在要取消		--> 更新分数和投票记录		差值的绝对值：1  +432
	2.之前投过赞成票，现在要取消		--> 更新分数和投票记录		差值的绝对值：1  -432
v=-1时，有两种情况
	1.之前没投过票，现在要投反对票		--> 更新分数和投票记录		差值的绝对值：1  -432
	2.之前投过赞成票，现在要改为反对票	--> 更新分数和投票记录		差值的绝对值：2  -432*2

投票的限制：
每个帖子子发表之日起一个星期之内允许用户投票，超过一个星期就不允许投票了
	1、到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2、到期之后删除那个 KeyPostVotedZSetPrefix
*/

func VoteForPost(userID int64, p *models.VoteDataForm) (err error) {
	// 1. 判断投票限制
	// 2. 更新帖子分数
	// 3. 记录用户为该帖子投票的数据
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int("Direction", p.Direction))
	return redis.VoteForPost(fmt.Sprintf("%d", userID), p.PostID, float64(p.Direction))
}
