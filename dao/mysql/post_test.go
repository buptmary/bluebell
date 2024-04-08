package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

func init() {
	dbCfg := settings.MySQLConfig{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "10428376",
		DB:       "bluebell-plus",
		Port:     3306,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		PostID:      10,
		AuthorID:    123,
		CommunityID: 2,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("Create post insert record into mysql failed, err:%v", err)
	}
	t.Logf("Create post insert record into mysql success")
}
