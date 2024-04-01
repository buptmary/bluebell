package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	// 1. 生成post_id(生成帖子ID)
	postID := snowflake.GenID()
	post.PostID = postID

	// 2. 创建帖子 保存到数据库
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(post)", zap.Error(err))
		return err
	}
	return err
}

// GetPostByID 根据Id查询帖子详情
func GetPostByID(postID int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合接口需要的数据
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(postID) failed", zap.Error(err))
		return
	}

	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed", zap.Int64("post.uid", post.AuthorID), zap.Error(err))
		return nil, err
	}

	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID(post.CommunityID) failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return nil, err
	}

	// 数据接口拼接
	data = &models.ApiPostDetail{ // (指针类型一定要做初始化)
		Post:            post,
		CommunityDetail: community,
		AuthorName:      user.UserName,
	}

	return data, nil
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) ([]*models.ApiPostDetail, error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed")
		return nil, err
	}
	// 初始化data返回值
	data := make([]*models.ApiPostDetail, 0, len(postList))

	// 向data填充信息
	for _, post := range postList {
		// 根据 authorID 获取作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed", zap.Int64("post.uid", post.AuthorID), zap.Error(err))
			return nil, err
		}
		// 根据 CommunityID 获取社区信息
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID(post.CommunityID) failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			return nil, err
		}

		// 数据接口拼接
		postDetail := &models.ApiPostDetail{
			Post:            post,
			CommunityDetail: community,
			AuthorName:      user.UserName,
		}
		data = append(data, postDetail)
	}
	return data, nil
}
