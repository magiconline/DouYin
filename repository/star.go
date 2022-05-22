package repository

import (
	"DouYin/logger"
	"sync"
)

type Star struct {
	ID      uint64
	UserId  uint64
	VideoId uint64
}

func (Star) TableName() string {
	return "star"
}

type StarDao struct {
}

var starDao *StarDao
var starOnce sync.Once

// NewStarDaoInstance 创建唯一实例
func NewStarDaoInstance() *StarDao {
	starOnce.Do(func() {
		starDao = &StarDao{}
	})
	return starDao
}

// AddStar 当进行点赞时，插入数据库。
func (*StarDao) AddStar(userId, videoId uint64) {
	star := Star{
		UserId:  userId,
		VideoId: videoId,
	}
	if err := DB.Table("star").Create(&star).Error; err != nil {
		logger.Logger.Printf("err", err)
		return
	}
}

// DeleteStar 当进行取消点赞时，删除数据库数据。
func (*StarDao) DeleteStar(userId, videoId uint64) {
	if err := DB.Table("star").Where("user_id = ? and video_id = ?", userId, videoId).Delete(&Star{}).Error; err != nil {
		logger.Logger.Printf("err", err)
		return
	}
}

// AuthorId 根据视频ID获取作者ID
func (*StarDao) AuthorId(videoID uint64) (*[]map[string]interface{}, error) {
	var authorID []map[string]interface{}
	result := DB.Table("video").Where("video_id = ?", videoID).Find(&authorID)
	return &authorID, result.Error
}

// AuthorInfo userId(user_id, user_name, follow_count, follower_count)字段
func (*StarDao) AuthorInfo(userId uint64) (*map[string]interface{}, error) {
	var author []map[string]interface{}
	result := DB.Table("users").Select("id", "user_name", "follow_count", "follower_count").Where("id = ?", userId).Find(&author)
	return &author[0], result.Error
}

// StarList StarVideoList 根据user_id查找用户点赞列表
func (*StarDao) StarList(userID uint64) (*[]map[string]interface{}, error) {
	var starList []map[string]interface{}
	result := DB.Table("star").Where("user_id = ?", userID).Find(&starList)
	return &starList, result.Error
}
