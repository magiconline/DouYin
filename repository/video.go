package repository

import (
	"fmt"
	// "DouYin/logger"
)

type VideoTable struct {
	VideoId       uint64 `gorm:"column:video_id"`
	UserId        uint64 `gorm:"column:user_id"`
	PlayUrl       string `gorm:"column:play_url"`
	CoverUrl      string `gorm:"column:cover_url"`
	FavoriteCount uint32 `gorm:"column:favorite_count"`
	CommentCount  uint32 `gorm:"column:comment_count"`
	UploadTime    uint64 `gorm:"column:upload_time"`
	Title         string `gorm:"column:title"`
}

func InsertVideoTable(videoTable *VideoTable) error {
	err := DB.Table("video").Create(&videoTable).Error
	if err != nil {
		fmt.Println("Insert VideoTable error:", err)
	}
	return err
}

func InsertVideo() error {
	return nil
}

func InsertCover() error {
	return nil
}

// 查找时间upload_time<latestTime, 降序排列的30条视频
func FeedAll(latestTime uint64) ([]map[string]interface{}, error) {
	var videoList []map[string]interface{}

	err := DB.Table("video").Where("upload_time < ?", latestTime).Order("upload_time desc").Limit(30).Find(&videoList).Error

	// fmt.Println(videoList)
	return videoList, err
}

// func FeedOne(latestTime uint32, userID uint64) {
// }

// 根据userID获取查询user表(user_id, user_name, follow_count, follower_count)字段
func AuthorInfo(userID uint64) (*map[string]interface{}, error) {
	var author []map[string]interface{}

	err := DB.Table("user").Select("user_id", "user_name", "follow_count", "follower_count", "title").Where("user_id = ?", userID).Find(&author).Error

	return &author[0], err
}

// 根据user_id查找所有视频
func UserVideoList(userID uint64) (*[]map[string]interface{}, error) {
	var videoList []map[string]interface{}
	err := DB.Table("video").Where("user_id = ?", userID).Order("upload_time desc").Find(&videoList).Error

	return &videoList, err
}
