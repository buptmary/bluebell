package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityByID(id int64) (communityDetail *models.CommunityDetail, err error) {
	return mysql.GetCommunityByID(id)
}

func GetCommunityList() (communityList []*models.Community, err error) {
	return mysql.GetCommunityList()
}
