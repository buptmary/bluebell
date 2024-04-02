package models

import "time"

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// Post 帖子Post结构体 内存对齐概念 字段类型相同的对齐 缩小变量所占内存大小
type Post struct {
	PostID      int64     `json:"post_id,string" gorm:"column:post_id"`
	AuthorID    int64     `json:"author_id,string"` // 序列化为string类型，防止数字失真
	CommunityID int64     `json:"community_id" gorm:"column:community_id" binding:"required"`
	Status      int32     `json:"status"`
	Title       string    `json:"title" gorm:"column:title" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateTime  time.Time `json:"update_time" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// ApiPostDetail 帖子返回详情结构体 拼接用户、社区和帖子信息
type ApiPostDetail struct {
	*Post                               // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息 json tag 将CommunityDetail 单独放入 community字段
	AuthorName       string             `json:"author_name"`
	VoteNum          int64              `json:"vote_num"` // 投票数量
}

// PostListForm 获取帖子列表query string参数
type PostListForm struct {
	Search      string `json:"search" form:"search"`
	CommunityID int64  `json:"community_id" form:"community_id"`
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
}
