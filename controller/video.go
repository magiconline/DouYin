package controller

import (
	"DouYin/service"
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
		"code":       0,
		"msg":        "success",
		"next_time":  nextTime,
		"video_list": videoList,
	}
}

func PublishAction(ctx *gin.Context) *interface{} {
	return nil
}

func PublishList(ctx *gin.Context) *interface{} {
	return nil
}
