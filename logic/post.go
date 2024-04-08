package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
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

	// 3. 向redis存储帖子信息
	if err := redis.CreatePost(post.PostID, post.CommunityID); err != nil {
		zap.L().Error("redis.CreatePost(post.PostID) failed", zap.Error(err))
		return err
	}

	return nil
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

// GetPostListInOrder 升级版帖子列表接口：按 创建时间 或者 分数排序
func GetPostListInOrder(p *models.PostListForm) (data []*models.ApiPostDetail, err error) {
	// 1. 根据参数中的排序规则去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		zap.L().Warn("post id list is 0")
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return nil, err
	}

	// 2. 根据id去MySQL数据库查询帖子详细信息
	// 返回的数据要按照给定的id顺序返回
	postList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//zap.L().Debug("ids", zap.Any("ids", ids))
	//zap.L().Debug("postlist", zap.Any("ids", postList))

	// 3. 查询每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Debug("redis.GetPostVoteData(ids)", zap.Any("ids", ids))
		return nil, err
	}

	// 3. 将帖子的作者和社区信息查询出来，填充到帖子中
	for idx, post := range postList {
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
			VoteNum:         voteData[idx],
		}
		data = append(data, postDetail)
	}
	return data, nil

}

func GetCommunityPostList(p *models.PostListForm) (data []*models.ApiPostDetail, err error) {
	// 1. 根据参数中的排序规则去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		zap.L().Warn("post id list is 0")
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder(p) return 0 data")
		return nil, err
	}
	zap.L().Debug("GetCommunityPostIDsInOrder", zap.Any("ids", ids))

	// 2. 根据id去MySQL数据库查询帖子详细信息
	// 返回的数据要按照给定的id顺序返回
	postList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 3. 查询每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Debug("redis.GetPostVoteData(ids)", zap.Any("ids", ids))
		return nil, err
	}

	data = make([]*models.ApiPostDetail, 0, len(postList))

	// 3. 将帖子的作者和社区信息查询出来，填充到帖子中
	for idx, post := range postList {
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
			VoteNum:         voteData[idx],
		}
		data = append(data, postDetail)
	}
	return data, nil
}

func GetPostListPro(p *models.PostListForm) (data []*models.ApiPostDetail, err error) {
	// 根据请求参数的不同，执行不同的业务逻辑
	// 如果CommunityID = 0，查询所有帖子列表
	if p.CommunityID == 0 {
		zap.L().Debug("Get post list by time/score")
		data, err = GetPostListInOrder(p)
	} else {
		// 根据社区id查询列表
		zap.L().Debug("get post list by community_id")
		data, err = GetCommunityPostList(p)
	}

	if err != nil {
		zap.L().Error("GetPostListPro failed", zap.Error(err))
		return nil, err
	}
	return data, nil
}
