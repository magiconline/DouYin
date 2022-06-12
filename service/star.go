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

//IsThumbUp 返回点赞状态
func IsThumbUp(userId, videoId uint64) (bool, error) {
	return repository.NewStarDaoInstance().IsThumbUp(userId, videoId)

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
		videoInfo, _ := repository.NewStarDaoInstance().VideoInfo(videoId)
		author, _ := AuthorInfo(videoInfo.UserId)
		//返回视频点赞状态 返回err 说明未点赞
		_, err1 := repository.NewStarDaoInstance().IsThumbUp(userID, videoId)
		if err1 != nil {
			continue
		}
		response_i := PublishActionResponse{
			ID:            videoInfo.VideoId,
			Author:        *author,
			PlayUrl:       server_ip + videoInfo.PlayUrl,
			CoverUrl:      server_ip + videoInfo.CoverUrl,
			FavoriteCount: videoInfo.FavoriteCount,
			CommentCount:  videoInfo.CommentCount,
			IsFavorite:    true,
			Title:         videoInfo.Title,
		}
		response = append(response, response_i)
	}
	return &response, nil
}
