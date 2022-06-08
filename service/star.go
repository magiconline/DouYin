package service

import (
	"DouYin/repository"
)

func AddStar(userId, videoId uint64) {
	repository.NewStarDaoInstance().AddStar(userId, videoId)
}

func DeleteStar(userId, videoId uint64) {
	repository.NewStarDaoInstance().DeleteStar(userId, videoId)
}

//StarVideoList 获取userID的所有的视频列表
func StarVideoList(token string, userID uint64) (*[]PublishActionResponse, error) {
	// 检查token
	if token != "" {
		_, err := Token2ID(token)
		if err != nil {
			return nil, err
		}
	}
	//根据userID获取点赞表
	var response []PublishActionResponse
	//根据用户ID获取视频ID列表
	starList, err := repository.NewStarDaoInstance().StarList(userID)
	if err != nil {
		return nil, err
	}
	for i := range *starList {
		//视频ID
		videoId := uint64((*starList)[i]["video_id"].(int64))
		//从视频ID获取视频信息
		videoInfo, err := repository.NewStarDaoInstance().VideoInfo(videoId)
		author, err := AuthorInfo(videoInfo.UserId)
		if err != nil {
			continue
		}
		//返回视频点赞状态
		stool, err := repository.NewStarDaoInstance().IsThumbUp(userID, videoId)
		if err != nil {
			continue
		}
		var isFavorite bool
		if stool == nil {
			isFavorite = false
		} else {
			isFavorite = true
		}
		response_i := PublishActionResponse{
			ID:            videoInfo.VideoId,
			Author:        *author,
			PlayUrl:       server_ip + videoInfo.PlayUrl,
			CoverUrl:      server_ip + videoInfo.CoverUrl,
			FavoriteCount: videoInfo.FavoriteCount,
			CommentCount:  videoInfo.CommentCount,
			IsFavorite:    isFavorite,
			Title:         videoInfo.Title,
		}
		response = append(response, response_i)
	}
	return &response, nil
}
