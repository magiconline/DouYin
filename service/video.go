package service

import (
	"DouYin/repository"
)

// 静态资源ip
// var server_ip = "http://172.26.59.240:8080"
var server_ip string

func SetServerIP(ip string) {
	server_ip = ip
}

type AuthorResponse struct {
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint64 `json:"follow_count"`
	FollowerCount uint64 `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type FeedResponse struct {
	ID            uint64         `json:"id"`
	Author        AuthorResponse `json:"author"`
	PlayUrl       string         `json:"play_url"`
	CoverUrl      string         `json:"cover_url"`
	FavoriteCount uint32         `json:"favorite_count"`
	CommentCount  uint32         `json:"comment_count"`
	IsFavorite    bool           `json:"is_favorite"`
	Title         string         `json:"title"`
}

type PublishActionResponse FeedResponse

func AuthorInfo(userID uint64) (*AuthorResponse, error) {
	author, err := repository.AuthorInfo(userID)
	if err != nil {
		return nil, err
	}
	return &AuthorResponse{ID: uint64(author.UserId), Name: author.UserName, FollowCount: 0, FollowerCount: 0, IsFollow: false}, nil
}

//Feed 获得视频流
// 如果token为空字符串则表示没有输入token，返回包含所有用户的视频流
// 如果token不为空，验证token，然后返回该用户的视频流
func Feed(latestTime uint64, token string) (uint64, *[]FeedResponse, error) {

	//获取当前用户, 验证token
	currentUserId, err := Token2ID(token)
	if err != nil {
		return 0, nil, err
	}
	var response []FeedResponse
	nextTime := latestTime // 如果没有新视频则nextTime = latestTime
	videoList, err := repository.FeedAll(latestTime)
	if err != nil {
		// 错误处理
		return 0, nil, err
	}
	// 将视频列表中填充author信息
	for i := range *videoList {
		userID := (*videoList)[i]["user_id"].(uint64)
		author, err := AuthorInfo(userID)
		if err != nil {
			continue
		}
		//返回视频点赞状态
		stool, _ := repository.NewStarDaoInstance().IsThumbUp(currentUserId, (*videoList)[i]["video_id"].(uint64))
		var isFavorite bool
		if stool == nil {
			isFavorite = false
		} else {
			isFavorite = true
		}
		response_i := FeedResponse{
			ID:            (*videoList)[i]["video_id"].(uint64),
			Author:        *author,
			PlayUrl:       server_ip + (*videoList)[i]["play_url"].(string),
			CoverUrl:      server_ip + (*videoList)[i]["cover_url"].(string),
			FavoriteCount: (*videoList)[i]["favorite_count"].(uint32),
			CommentCount:  (*videoList)[i]["comment_count"].(uint32),
			IsFavorite:    isFavorite,
			Title:         (*videoList)[i]["title"].(string),
		}
		response = append(response, response_i)
		nextTime = (*videoList)[i]["upload_time"].(uint64)
	}
	return nextTime, &response, nil
}

// UserVideoList 获取userID的所有的视频列表
func UserVideoList(token string, userID uint64) (*[]PublishActionResponse, error) {
	// 检查token
	if token != "" {
		_, err := Token2ID(token)
		if err != nil {
			return nil, err
		}
	}

	author, err := AuthorInfo(userID)
	if err != nil {
		return nil, err
	}

	// 根据userID查找数据库，按上传时间排序
	videoList, err := repository.UserVideoList(userID)
	if err != nil {
		return nil, err
	}

	var response []PublishActionResponse
	// 将视频列表中填充author信息
	for i := range *videoList {
		response_i := PublishActionResponse{
			ID:            (*videoList)[i]["video_id"].(uint64),
			Author:        *author,
			PlayUrl:       server_ip + (*videoList)[i]["play_url"].(string),
			CoverUrl:      server_ip + (*videoList)[i]["cover_url"].(string),
			FavoriteCount: (*videoList)[i]["favorite_count"].(uint32),
			CommentCount:  (*videoList)[i]["comment_count"].(uint32),
			IsFavorite:    false,
			Title:         (*videoList)[i]["title"].(string),
		}
		response = append(response, response_i)
	}

	return &response, nil
}
