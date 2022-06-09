package controller

import (
	"DouYin/repository"
	"DouYin/service"
	"os"
	"strconv"
	"strings"
	"time"

	"DouYin/logger"

	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
)

//Feed 获得视频流，获得参数，调用service层，返回map
func Feed(ctx *gin.Context) *gin.H {
	// 获得参数
	latestTimeStr := ctx.Query("latest_time")
	token := ctx.Query("token")

	// 类型转换
	latestTimeInt, err := strconv.Atoi(latestTimeStr)
	if err != nil {
		logger.Logger.Println("error:", err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}
	nextTime, videoList, err := service.Feed(uint64(latestTimeInt), token)
	// 获得视频流错误
	if err != nil {
		logger.Logger.Println("error:", err)
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}
	// 正常返回
	return &gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"next_time":   nextTime,
		"video_list":  *videoList,
	}
}

//PublishAction 投稿接口
func PublishAction(ctx *gin.Context) *gin.H {
	data, err := ctx.FormFile("data")
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}
	token := ctx.PostForm("token")
	title := ctx.PostForm("title")

	// 通过文件后缀名验证格式
	if arr := strings.Split(data.Filename, "."); arr[len(arr)-1] != "mp4" {
		return &gin.H{
			"status_code": 1,
			"status_msg":  "视频格式不支持",
		}

	}

	// 验证token
	userID, err := service.Token2ID(token)
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	// 生成保存路径
	curTime := time.Now()
	path := "/static/" + curTime.Format("2006/01/02") + "/"
	name := uuid.NewString()
	videoName := name + ".mp4"
	coverName := name + ".jpg"

	err = os.MkdirAll("."+path, 0777)
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	// 保存视频
	err = ctx.SaveUploadedFile(data, "."+path+videoName)
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	// 生成并保存封面
	if err = repository.InsertCover(path, videoName, coverName); err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	// 插入数据库
	videoTable := repository.VideoTable{
		UserId:     userID,
		PlayUrl:    path + videoName,
		CoverUrl:   path + coverName,
		UploadTime: uint64(curTime.UnixMilli()),
		Title:      title,
	}

	err = repository.InsertVideoTable(&videoTable)
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	return &gin.H{
		"status_code": 0,
		"status_msg":  "success",
	}
}

// PublishList 登录用户的视频发布列表，直接列出用户所有投稿过的视频
func PublishList(ctx *gin.Context) *gin.H {
	// 获得参数
	token := ctx.Query("token")
	userIDStr := ctx.Query("user_id")

	// 类型转换
	userID, err := strconv.ParseUint(userIDStr, 10, 0)
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}
	// 获得信息
	videoList, err := service.UserVideoList(token, userID)
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	return &gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"video_list":  videoList,
	}
}
