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
