package controller

import (
	"DouYin/logger"
	"DouYin/service"
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

	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {
		logger.Logger.Println("error:", err)
		return &gin.H{
			"status_code": 1,
			"status_msg":  "用户ID格式错误",
		}
	}
	videoIdInt, err1 := strconv.Atoi(videoIdStr)
	if err1 != nil {
		return &gin.H{
			"status_code": 2,
			"status_msg":  "视频ID格式错误",
		}
	}
	actionTypeInt, err2 := strconv.Atoi(actionTypeStr)
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
