package service

import (
	"DouYin/repository"
	"fmt"
)

func AddStar(userId, videoId uint64) {
	repository.NewStarDaoInstance().AddStar(userId, videoId)
}

func DeleteStar(userId, videoId uint64) {
	repository.NewStarDaoInstance().DeleteStar(userId, videoId)
}

// StarVideoList 获取点赞视频列表
func StarVideoList(userId uint64) (*[]map[string]interface{}, error) {
	//获取点赞的视频Id列表
	starVideoList, err := repository.NewStarDaoInstance().StarList(userId)
	var authorId *[]map[string]interface{}
	fmt.Println(starVideoList)
	if err != nil {
		// 错误处理，返回空列表
		return starVideoList, err
	}
	// 将视频列表中填充author信息
	for i := range *starVideoList {
		videoId := uint64((*starVideoList)[i]["video_id"].(int64))
		authorId, err = repository.NewStarDaoInstance().AuthorId(videoId)
		authorIdInfo := (*authorId)[0]["user_id"].(uint64)
		var author *map[string]interface{}
		author, err = repository.NewStarDaoInstance().AuthorInfo(authorIdInfo)
		fmt.Println("author：", author)
		if err != nil {
			// 错误处理
		}
		(*author)["id"] = (*author)["id"]
		(*author)["name"] = (*author)["user_name"]
		(*author)["is_follow"] = false
		delete(*author, "user_name")
		(*authorId)[i]["author"] = author
		(*authorId)[i]["id"] = videoId
		// 将相对url转为完整url
		(*authorId)[i]["play_url"] = server_ip + (*authorId)[i]["play_url"].(string)
		(*authorId)[i]["cover_url"] = server_ip + (*authorId)[i]["cover_url"].(string)
		(*authorId)[i]["favorite_count"] = (*authorId)[i]["favorite_count"].(uint32)
		(*authorId)[i]["is_favorite"] = true
		(*authorId)[i]["title"] = "测试主题"
		delete((*authorId)[i], "video_id")
		delete((*authorId)[i], "upload_time")
	}
	return authorId, nil
}
