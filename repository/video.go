package repository

import (
	"DouYin/logger"
	"bytes"
	"fmt"
	"image"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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
		logger.Println(err.Error())
	}
	return err
}

// 生成并保存图像
func InsertCover(path string, videoName string, coverName string) error {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input("."+path+videoName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 0)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run()

	if err != nil {
		return err
	}

	var image image.Image
	image, err = imaging.Decode(buf)
	if err != nil {
		return err
	}

	err = imaging.Save(image, "."+path+coverName)
	if err != nil {
		return err
	}

	return nil
}

// FeedAll 查找时间upload_time<latestTime, 降序排列的30条视频
func FeedAll(latestTime uint64) (*[]map[string]interface{}, error) {
	var videoList []map[string]interface{}

	err := DB.Table("video").Where("upload_time < ?", latestTime).Order("upload_time desc").Limit(30).Find(&videoList).Error

	return &videoList, err
}

// AuthorInfo 根据userID获取查询user表(user_id, user_name, follow_count, follower_count)字段
func AuthorInfo(userID uint64) (*User, error) {
	var user User

	err := DB.Table("user").Select("user_id", "user_name", "follow_count", "follower_count").Where("user_id = ?", userID).First(&user).Error

	return &user, err
}

// UserVideoList 根据user_id查找所有视频
func UserVideoList(userID uint64) (*[]map[string]interface{}, error) {
	var videoList []map[string]interface{}
	err := DB.Table("video").Where("user_id = ?", userID).Order("upload_time desc").Find(&videoList).Error
	return &videoList, err
}
