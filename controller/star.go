package controller

import (
	"DouYin/logger"
	"DouYin/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Favorite(ctx *gin.Context) *gin.H {
	videoIdStr := ctx.Query("video_id")
	token := ctx.Query("token")
/*	k := videoIdStr
	v := "null"
	if err := repository.GetRedisLock(k, v, 10*time.Second); err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}
	defer repository.DeleteRedisLock(k, v)*/
	//1.点赞 2.取消点赞
	actionTypeStr := ctx.Query("action_type")
	if token == "" {
		return &gin.H{
			"status_code": 1,
			"status_msg":  "用户未登录",
		}
	}
	_, err := service.ParseToken(token)
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err,
		}
	}
	currentUserID, err := service.Token2ID(token)
	if err != nil {
		logger.Logger.Println("error:", err)
		return &gin.H{
			"status_code": 1,
			"status_msg":  err,
		}
	}
	videoIdInt, err1 := strconv.ParseUint(videoIdStr, 10, 64)
	if err1 != nil {
		return &gin.H{
			"status_code": 2,
			"status_msg":  err1,
		}
	}
	actionTypeInt, err2 := strconv.ParseUint(actionTypeStr, 10, 8)
	if err2 != nil {
		return &gin.H{
			"status_code": 3,
			"status_msg":  err2,
		}
	}
	if actionTypeInt != 1 && actionTypeInt != 2 {
		return &gin.H{
			"status_code": 4,
			"status_msg":  "传入参数为1或者2！",
		}
	}
	if actionTypeInt == 1 {
		//查询数据库，获取点赞状态
		flag, err := service.IsThumbUp(currentUserID, videoIdInt)
		if err != nil {
			return &gin.H{
				"status_code": 1,
				"status_msg":  err,
			}
		}
		if flag {
			return &gin.H{
				"status_code": 5,
				"status_msg":  "请勿重复点赞！",
			}
		}
		service.AddStar(currentUserID, videoIdInt)
		return &gin.H{
			"status_code": 0,
			"status_msg":  "点赞成功！",
		}
	} else {
		//查询数据库，获取点赞状态
		flag, err := service.IsThumbUp(currentUserID, videoIdInt)
		if err != nil {
			return &gin.H{
				"status_code": 1,
				"status_msg":  err.Error(),
			}
		}
		if !flag {
			return &gin.H{
				"status_code": 6,
				"status_msg":  "当前暂无点赞数据！",
			}
		}
		service.DeleteStar(currentUserID, videoIdInt)
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
	// fmt.Println("token：", token)
	if token == "" {
		return &gin.H{
			"status_code": 1,
			"status_msg":  "用户未登录",
		}
	}
	// 类型转换,查看用户ID
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return &gin.H{
			"status_code": 2,
			"status_msg":  err.Error(),
		}
	}
	// 获得信息
	videoList, err := service.StarVideoList(token, userID)
	if err != nil {
		// 错误处理
		return &gin.H{
			"status_code": 3,
			"status_msg":  err.Error(),
		}
	}
	return &gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"video_list":  videoList,
	}
}
