package models

import "time"

type Community struct {
	CommunityID   int64  `json:"community_id" gorm:"column:community_id"`
	CommunityName string `json:"community_name" gorm:"column:community_name"`
}

type CommunityDetail struct {
	CommunityID   int64     `json:"community_id" gorm:"column:community_id"`
	CommunityName string    `json:"community_name" gorm:"column:community_name"`
	Introduction  string    `json:"introduction,omitempty" gorm:"column:introduction"`
	CreateTime    time.Time `json:"create_time" gorm:"column:create_time"`
}
