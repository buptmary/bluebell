package mysql

import (
	"bluebell/models"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetCommunityList  返回社区ID、名称列表
func GetCommunityList() (communityList []*models.Community, err error) {
	// 初始化slice
	communityList = make([]*models.Community, 0)
	err = db.Find(&communityList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Warn("there is no community!")
			return nil, nil
		}
		return nil, err
	}
	return communityList, err
}

// GetCommunityByID 根据社区ID获取社区详情
func GetCommunityByID(id int64) (communityDetail *models.CommunityDetail, err error) {
	// 参数指针变量需要初始化
	communityDetail = new(models.CommunityDetail)
	err = db.Table("communities").Where("community_id = ?", id).First(communityDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("record not found")
			return nil, ErrorInvalidID
		}
		zap.L().Error("get community by id failed", zap.Error(err))
		return nil, ErrorQueryFailed
	}
	return communityDetail, nil
}
