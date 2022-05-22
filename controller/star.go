package controller

import (
	"DouYin/logger"
	"DouYin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Favorite(ctx *gin.Context) *gin.H {
	userIdStr := ctx.Query("user_id")
	videoIdStr := ctx.Query("video_id")
	//获取token
	token := ctx.Query("token")
	//1.点赞 2.取消点赞
	actionTypeStr := ctx.Query("action_type")

	userIdInt, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		logger.Logger.Println("error:", err)
		return &gin.H{
			"status_code": 1,
			"status_msg":  "用户ID格式错误",
		}
	}
	videoIdInt, err1 := strconv.ParseUint(videoIdStr, 10, 64)
	if err1 != nil {
		return &gin.H{
			"status_code": 2,
			"status_msg":  "视频ID格式错误",
		}
	}
	actionTypeInt, err2 := strconv.ParseUint(actionTypeStr, 10, 8)
	if err2 != nil {
		return &gin.H{
			"status_code": 3,
			"status_msg":  "点赞格式传入错误",
		}
	}
	if actionTypeInt != 1 && actionTypeInt != 2 {
		return &gin.H{
			"status_code": 4,
			"status_msg":  "传入参数为1或者2！",
		}
	}
	logger.Logger.Printf("userIdInt：%d actionTypeInt：%d videoIdInt：%d token：%s", userIdInt, actionTypeInt, videoIdInt, token)
	//用户鉴权token
	if actionTypeInt == 1 {
		service.AddStar(uint64(userIdInt), uint64(videoIdInt))
		return &gin.H{
			"status_code": 0,
			"status_msg":  "点赞成功！",
		}
	} else {
		service.DeleteStar(uint64(userIdInt), uint64(videoIdInt))
		return &gin.H{
			"status_code": 0,
			"status_msg":  "取消点赞成功！",
		}
	}
}

// ThumbListVideo 视频列表
func ThumbListVideo(ctx *gin.Context) *gin.H {
	// 获得参数
	token := ctx.Query("token")
	userIDStr := ctx.Query("user_id")
	fmt.Println("token：", token)
	// 类型转换
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  "failed",
		}
	}
	// 获得信息
	videoList, err := service.StarVideoList(userID)
	if err != nil {
		// 错误处理
		return &gin.H{
			"status_code": 1,
			"status_msg":  "failed",
		}
	}
	return &gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"video_list":  videoList,
	}
}
