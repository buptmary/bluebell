package mysql

import (
	"bluebell/models"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	err = db.Create(post).Error
	if err != nil {
		zap.L().Error("mysql create post failed", zap.Error(err))
		return err
	}
	return nil
}

// GetPostByID 根据Id查询帖子详情
func GetPostByID(postID int64) (post *models.Post, err error) {
	post = new(models.Post)
	err = db.Where("post_id = ?", postID).First(post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorInvalidID
		}
		zap.L().Error("mysql get post by id failed", zap.Error(err))
		return nil, ErrorQueryFailed
	}
	return post, nil
}

// GetPostList 获取帖子列表-分页
func GetPostList(page, size int64) ([]*models.Post, error) {
	postList := make([]*models.Post, size)
	// 计算offset，用于分页查询
	offset := int((page - 1) * size)
	limit := int(size)

	err := db.Offset(offset).Limit(limit).Find(&postList).Error
	if err != nil {
		zap.L().Error("mysql get post list failed", zap.Error(err))
		return nil, err
	}
	return postList, nil
}
