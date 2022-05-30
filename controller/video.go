package controller

import (
	"DouYin/service"
	"fmt"
	"strconv"

	"DouYin/logger"

	"github.com/gin-gonic/gin"
)

// 获得视频流，获得参数，调用service层，返回map
func Feed(ctx *gin.Context) *gin.H {
	// 获得参数
	latestTimeStr := ctx.Query("latest_time")
	token := ctx.DefaultQuery("token", "")

	// 类型转换
	// latestTimeUint32, err := strconv.ParseUint(latestTimeStr, 10, 32)
	latestTimeInt, err := strconv.Atoi(latestTimeStr)
	if err != nil {
		logger.Logger.Println("error:", err)
	}
	nextTime, videoList, err := service.Feed(uint64(latestTimeInt), token)

	// 获得视频流错误
	if err != nil {
		logger.Logger.Println("error:", err)
		return &gin.H{
			"code": 1,
			"msg":  "failed",
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

// 投稿接口
func PublishAction(ctx *gin.Context) *gin.H {
	data, err := ctx.FormFile("data")
	if err != nil {
		fmt.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}
	token := ctx.PostForm("token")
	title := ctx.PostForm("title")

	err = service.PublishAction(data, token, title)
	if err != nil {
		fmt.Println(err.Error())
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

// 登录用户的视频发布列表，直接列出用户所有投稿过的视频
func PublishList(ctx *gin.Context) *gin.H {
	// 获得参数
	token := ctx.Query("token")
	userIDStr := ctx.Query("user_id")

	// 类型转换
	userID, _ := strconv.ParseUint(userIDStr, 10, 0)

	// 获得信息
	videoList, err := service.UserVideoList(token, userID)
	if err != nil {
		// 错误处理

		return &gin.H{
			"code": 1,
			"msg":  "success",
		}
	}

	return &gin.H{
		"code":       0,
		"msg":        "failed",
		"video_list": videoList,
	}
}
