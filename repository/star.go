package repository

import (
	"DouYin/logger"
	"gorm.io/gorm"
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

// AddStar 当进行点赞时，插入数据库,同时更新点赞数
func (*StarDao) AddStar(userId, videoId uint64) {
	star := Star{
		UserId:  userId,
		VideoId: videoId,
	}
	//开启事务
	DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("star").Create(&star).Error; err != nil {
			logger.Logger.Printf("err", err)
			return err
		}
		var count int64
		if err := tx.Table("star").Where("video_id = ?", videoId).Count(&count).Error; err != nil {
			logger.Logger.Printf("err", err)
			return err
		}
		if err := tx.Table("video").Where("video_id = ?", videoId).Update("favorite_count", count).Error; err != nil {
			logger.Logger.Printf("err", err)
			return err
		}
		return nil
	})
}

// DeleteStar 当进行取消点赞时，删除数据库数据。
func (*StarDao) DeleteStar(userId, videoId uint64) {
	//开启事务
	DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("star").Where("user_id = ? and video_id = ?", userId, videoId).Delete(&Star{}).Error; err != nil {
			logger.Logger.Printf("err", err)
			return err
		}
		var count int64
		if err := tx.Table("star").Where("video_id = ?", videoId).Count(&count).Error; err != nil {
			logger.Logger.Printf("err", err)
			return err
		}
		if err := tx.Table("video").Where("video_id = ?", videoId).Update("favorite_count", count).Error; err != nil {
			logger.Logger.Printf("err", err)
			return err
		}
		return nil
	})
}

// VideoInfo 根据视频ID获取视频信息
func (*StarDao) VideoInfo(videoID uint64) (*VideoTable, error) {
	var video VideoTable
	result := DB.Table("video").Where("video_id = ?", videoID).Find(&video)
	return &video, result.Error
}

// AuthorInfo userId(user_id, user_name, follow_count, follower_count)字段
func (*StarDao) AuthorInfo(userId uint64) (*User, error) {
	var user User
	err := DB.Table("user").Select("user_id", "user_name", "follow_count", "follower_count").Where("user_id = ?", userId).First(&user).Error
	return &user, err
}

// StarList StarVideoList 根据user_id查找用户点赞列表
func (*StarDao) StarList(userID uint64) (*[]map[string]interface{}, error) {
	var starList []map[string]interface{}
	result := DB.Table("star").Where("user_id = ?", userID).Find(&starList)
	return &starList, result.Error
}

//IsThumbUp 返回点赞状态
func (*StarDao) IsThumbUp(userID, videoID uint64) (*Star, error) {
	var star Star
	err := DB.Table("star").Where("user_id = ? and video_id = ?", userID, videoID).Find(&star).Error
	return &star, err
}

/*//FavoriteCount 返回视频获赞总数
func (*StarDao) FavoriteCount(videoID uint64) (int64, error) {
	var count int64
	DB.Model(&Star{}).Where("video_id = ?", videoID).Count(&count)
	return count, nil
}*/

////UpdateFavoriteCount 更新视频获赞总数
//func (*StarDao) UpdateFavoriteCount(videoID, favoriteCount uint64) {
//	DB.Model(&VideoTable{}).Where("video_id = ?", videoID).Update("favorite_count", favoriteCount)
//}
