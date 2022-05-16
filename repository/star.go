package repository

import (
	"sync"
)

type Star struct {
	ID      uint64
	UserId  uint64
	VideoId uint8
	Status  uint8
}

func (Star) TableName() string {
	return "star"
}

type StarDao struct {
}

var starDao *StarDao
var starOnce sync.Once

//创建唯一实例
func NewStarDaoInstance() *StarDao {
	starOnce.Do(func() {
		starDao = &StarDao{}
	})
	return starDao
}

func (*StarDao) FavorableOperation(star *Star) (error, bool) {
	err := DB.Model(star).Update("status", star.Status).Where(&Star{UserId: star.UserId, VideoId: star.VideoId})
	if err != nil {
		//输出错误信息
		return nil, false
	}
	return nil, true
}
func (*StarDao) QueryByUserIdAndVideoId(userId, videoId uint64) *Star {
	var star Star
	DB.Where("user_id = ? and video_id = ?", userId, videoId).Find(&star)
	return &star
}
