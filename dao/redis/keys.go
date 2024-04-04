package redis

// redis key
// 使用命名空间的方式，方便查询与拆分
// 和其他项目共用Redis时，通过项目名进行区分
const (
	KeyPostTimeZSet        = "bluebell:post:time"   // 帖子及发帖时间
	KeyPostScoreZSet       = "bluebell:post:score"  // 帖子及投票分数
	KeyPostVotedZSetPrefix = "bluebell:post:voted:" // 记录用户及投票类型;参数是post_id
	KeyCommunitySetPrefix  = "bluebell:community:"  // Set保存每个分区下帖子的id
)
