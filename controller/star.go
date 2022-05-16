package controller

import (
	"DouYin/repository"
	"DouYin/service"
	"strconv"
)

type FavoriteData struct {
	StatusCode uint
	StatusMsg  string
}

func Favorite(userId, token, videoId, actionType string) *FavoriteData {
	user_id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return &FavoriteData{
			StatusCode: 1,
			StatusMsg:  "用户ID格式错误",
		}
	}
	video_id, err := strconv.ParseUint(videoId, 10, 64)
	if err != nil {
		return &FavoriteData{
			StatusCode: 2,
			StatusMsg:  "视频ID格式错误",
		}
	}
	//action_type, err := strconv.ParseUint(actionType, 10, 8)
	if err != nil {
		return &FavoriteData{
			StatusCode: 3,
			StatusMsg:  "点赞状态格式错误",
		}
	}
	//用户鉴权token
	//查询当前用户此视频的点赞状态
	var current_star *repository.Star
	current_star = service.QueryByUserIdAndVideoId(user_id, video_id)
	//根据用户ID及视频ID修改点赞状态
	if current_star == nil {
		//此时用户未曾点赞 插入数据
	}
	//service.FavoriteOperation()
	return nil
}
