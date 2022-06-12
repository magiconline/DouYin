package repository

import (
	"DouYin/logger"
	"sync"

	"gorm.io/gorm"
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
func (*StarDao) AddStar(userId, videoId uint64) error {
	star := Star{
		UserId:  userId,
		VideoId: videoId,
	}
	//开启事务
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	//插入star表数据
	if err := tx.Table("star").Create(&star).Error; err != nil {
		logger.Logger.Printf("err, %v", err)
		return err
	}
	//锁住指定video_id的记录
	video := VideoTable{}
	if err := tx.Table("video").Set("gorm:query_option", "FOR UPDATE").First(&video, videoId).Error; err != nil {
		tx.Rollback()
		return err
	}
	//更新操作
	tx.Table("video").Where("video_id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
	// commit事务，释放锁
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// DeleteStar 当进行取消点赞时，删除数据库数据。
func (*StarDao) DeleteStar(userId, videoId uint64) error {
	//开启事务
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	//删除star表数据
	if err := tx.Table("star").Where("user_id = ? and video_id = ?", userId, videoId).Delete(&Star{}).Error; err != nil {
		logger.Logger.Printf("err", err)
		return err
	}
	//锁住指定video_id的记录
	video := VideoTable{}
	if err := tx.Table("video").Set("gorm:query_option", "FOR UPDATE").First(&video, videoId).Error; err != nil {
		tx.Rollback()
		return err
	}
	//更新操作
	tx.Table("video").Where("video_id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
	// commit事务，释放锁
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
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
func (*StarDao) IsThumbUp(userID, videoID uint64) (bool, error) {
	var star Star
	result := DB.Table("star").Where("user_id = ? and video_id = ?", userID, videoID).Limit(1).Find(&star)

	return result.RowsAffected != 0, result.Error
}
