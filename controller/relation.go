package controller

import (
	"DouYin/logger"
	"DouYin/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 关注操作
func RelationAction(ctx *gin.Context) *gin.H {
	// 获取参数
	token := ctx.Query("token")
	toUserIDStr := ctx.Query("to_user_id")
	actionTypeStr := ctx.Query("action_type")

	toUserID, err := strconv.ParseUint(toUserIDStr, 10, 64)
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	// 关注操作
	err = service.RelationAction(token, toUserID, actionTypeStr == "1")
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	return &gin.H{
		"status_code": 0,
		"status_msg":  "ok",
	}
}

// 获取关注列表
func FollowList(ctx *gin.Context) *gin.H {
	// 获取参数
	userIDStr := ctx.Query("user_id")
	token := ctx.Query("token")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	// 验证token
	_, err = service.Token2ID(token)
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	// 获取关注列表
	followList, err := service.FollowList(userID)
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}
	return &gin.H{
		"status_code": 0,
		"statud_msg":  "ok",
		"user_list":   followList,
	}
}

// 获取粉丝列表
func FollowerList(ctx *gin.Context) *gin.H {
	// 获取参数
	userIDStr := ctx.Query("user_id")
	token := ctx.Query("token")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	// 验证token
	_, err = service.Token2ID(token)
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	// 获取粉丝列表
	followerList, err := service.FollowerList(userID)
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}
	return &gin.H{
		"status_code": 0,
		"statud_msg":  "ok",
		"user_list":   followerList,
	}
}
