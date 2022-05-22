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
}

func InsertVideoTable(videoTable *VideoTable) error {
	result := DB.Table("video").Create(&videoTable)
	if result.Error != nil {

		fmt.Println("InsertVideoTable error:", result.Error)
		return result.Error
	}
	return nil
}

func InsertVideo() error {
	return nil
}

func InsertCover() error {
	return nil
}

// FeedAll 查找时间upload_time<latestTime, 降序排列的30条视频
func FeedAll(latestTime uint64) ([]map[string]interface{}, error) {
	var videoList []map[string]interface{}

	result := DB.Table("video").Where("upload_time < ?", latestTime).Order("upload_time desc").Limit(30).Find(&videoList)

	// fmt.Println(videoList)
	return videoList, result.Error
}

// func FeedOne(latestTime uint32, userID uint64) {
// }

// AuthorInfo 根据userID获取查询user表(user_id, user_name, follow_count, follower_count)字段
func AuthorInfo(userID uint64) (*map[string]interface{}, error) {
	var author []map[string]interface{}

	result := DB.Table("user").Select("user_id", "user_name", "follow_count", "follower_count").Where("user_id = ?", userID).Find(&author)

	return &author[0], result.Error
}
