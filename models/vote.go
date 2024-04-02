package models

// VoteDataForm 投票数据
type VoteDataForm struct {
	//UserID int64 从请求中获取当前用户id
	PostID    string `json:"post_id" binding:"required"`              // 帖子id
	Direction int    `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票(1) 反对票(-1) 取消投票(0)
}
