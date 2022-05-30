package repository

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"

	// "DouYin/logger"
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
		fmt.Println("Insert VideoTable error:", err)
	}
	return err
}

// 插入视频
func InsertVideo(path string, videoName string, video *multipart.FileHeader) error {
	// 打开视频句柄
	file, err := video.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	// 检测文件夹是否创建
	err = os.MkdirAll("."+path, 0777)
	if err != nil {
		return err
	}

	// 本地创建文件，如果文件已存在则会被清空
	localFile, err := os.Create("." + path + videoName)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// 拷贝文件
	_, err = io.Copy(localFile, file)
	return err
}

// 生成并保存图像
func InsertCover(path string, videoName string, coverName string) error {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input("."+path+videoName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 0)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
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

// func FeedOne(latestTime uint32, userID uint64) {
// }

// AuthorInfo 根据userID获取查询user表(user_id, user_name, follow_count, follower_count)字段
func AuthorInfo(userID uint64) (*map[string]interface{}, error) {
	var author []map[string]interface{}

	err := DB.Table("user").Select("user_id", "user_name", "follow_count", "follower_count", "title").Where("user_id = ?", userID).Find(&author).Error

	if err != nil {
		return nil, err
	}

	return &author[0], nil
}

// 根据user_id查找所有视频
func UserVideoList(userID uint64) (*[]map[string]interface{}, error) {
	var videoList []map[string]interface{}
	err := DB.Table("video").Where("user_id = ?", userID).Order("upload_time desc").Find(&videoList).Error

	return &videoList, err
}
