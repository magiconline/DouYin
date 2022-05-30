package service

import (
	"DouYin/repository"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	uuid "github.com/google/uuid"
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

func AuthorInfo(userID uint64) (*AuthorResponse, error) {
	author, err := repository.AuthorInfo(userID)
	if err != nil {
		return nil, err
	}
	return &AuthorResponse{ID: uint64(author.UserId), Name: author.UserName, FollowCount: 0, FollowerCount: 0, IsFollow: false}, nil
}

// 获得视频流
// 如果token为空字符串则表示没有输入token，返回包含所有用户的视频流
// 如果token不为空，验证token，然后返回该用户的视频流
func Feed(latestTime uint64, token string) (uint64, *[]FeedResponse, error) {
	if token == "" {
		var response []FeedResponse
		nextTime := latestTime
		// 获得视频列表
		videoList, err := repository.FeedAll(latestTime)
		if err != nil {
			// 错误处理
			return latestTime, nil, err
		}

		// 将视频列表中填充author信息
		for i := range videoList {
			userID := videoList[i]["user_id"].(uint64)

			author, err := AuthorInfo(userID)
			if err != nil {
				fmt.Println("获取AuthorInfo错误:", err)
				continue
			}
			response_i := FeedResponse{
				ID:            userID,
				Author:        *author,
				PlayUrl:       server_ip + videoList[i]["play_url"].(string),
				CoverUrl:      server_ip + videoList[i]["cover_url"].(string),
				FavoriteCount: videoList[i]["favorite_count"].(uint32),
				CommentCount:  videoList[i]["comment_count"].(uint32),
				IsFavorite:    false,
				Title:         videoList[i]["title"].(string),
			}
			response = append(response, response_i)
			nextTime = videoList[i]["upload_time"].(uint64)
			fmt.Println("nextTime:", nextTime)
		}

		return nextTime, &response, nil
	} else {
		return latestTime, nil, errors.New("暂不支持的方法！")
	}
}

// 获取userID的所有的视频列表
func UserVideoList(token string, userID uint64) (*[]map[string]interface{}, error) {
	// 检查token
	_, err := Token2ID(token)
	if err != nil {
		fmt.Println("token验证失败:", err)
		return nil, err
	}

	// 匹配, 根据userID查找数据库，按上传时间排序
	videoList, err := repository.UserVideoList(userID)
	if err != nil {
		// 错误处理，返回空列表
		return videoList, err
	}
	// 将视频列表中填充author信息
	for i := range *videoList {

		user_id := (*videoList)[i]["user_id"].(uint64)

		author, err := repository.AuthorInfo(user_id)
		if err != nil {
			// 错误处理
			fmt.Println(err)
			continue
		}
		authorMap := make(map[string]interface{})
		authorMap["id"] = author.UserId
		authorMap["name"] = author.UserName
		authorMap["is_follow"] = false

		(*videoList)[i]["author"] = authorMap

		// 将相对url转为完整url
		(*videoList)[i]["play_url"] = server_ip + (*videoList)[i]["play_url"].(string)
		(*videoList)[i]["cover_url"] = server_ip + (*videoList)[i]["cover_url"].(string)
	}

	return videoList, nil
	// 不匹配
}

// 登录用户选择视频上传
func PublishAction(data *multipart.FileHeader, token string, title string) error {
	// 验证token
	userID, err := Token2ID(token)
	if err != nil {
		return err
	}

	// 生成保存路径
	path := "/static/" + time.Now().Format("2006/01/02") + "/"
	name := uuid.NewString()
	videoName := name + ".mp4"
	coverName := name + ".jpg"

	// 保存视频
	err = repository.InsertVideo(path, videoName, data)
	if err != nil {
		return err
	}

	// 生成并保存缩略图
	err = repository.InsertCover(path, videoName, coverName)
	if err != nil {
		return err
	}

	// 插入数据库
	videoTable := repository.VideoTable{
		UserId:     userID,
		PlayUrl:    path + videoName,
		CoverUrl:   path + coverName,
		UploadTime: uint64(time.Now().UnixMilli()),
		Title:      title,
	}

	err = repository.InsertVideoTable(&videoTable)
	return err
}
