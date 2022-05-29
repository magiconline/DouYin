package service

import (
	"DouYin/repository"
)

func InsertVideo() error {
	return nil
}

// 静态资源ip
// var server_ip = "http://172.26.59.240:8080"
var server_ip string

func SetServerIP(ip string) {
	server_ip = ip
}

// 获得视频流
// 如果token为空字符串则表示没有输入token，返回包含所有用户的视频流
// 如果token不为空，验证token，然后返回该用户的视频流
func Feed(latestTime uint64, token string) (uint64, []map[string]interface{}, error) {
	if token == "" {
		// 获得视频列表
		videoList, err := repository.FeedAll(latestTime)
		if err != nil {
			// 错误处理
			return 0, nil, err
		}

		// 将视频列表中填充author信息
		for i := range videoList {

			user_id := videoList[i]["user_id"].(uint64)

			var author *map[string]interface{}
			author, err = repository.AuthorInfo(user_id)
			if err != nil {
				// 错误处理
			}
			(*author)["id"] = (*author)["user_id"]
			(*author)["name"] = (*author)["user_name"]
			(*author)["is_follow"] = false

			videoList[i]["author"] = author

			// 将相对url转为完整url
			videoList[i]["play_url"] = server_ip + videoList[i]["play_url"].(string)
			videoList[i]["cover_url"] = server_ip + videoList[i]["cover_url"].(string)
		}

		// 获得nextTime
		nextTime := videoList[len(videoList)-1]["upload_time"].(uint64)

		return nextTime, videoList, nil
	} else {
		// 验证token

		// 验证成功

		// 验证失败
		return 0, nil, nil
	}
}

// 根据userID和token获得所有的视频列表
func UserVideoList(token string, userID uint64) (*[]map[string]interface{}, error) {
	// 检查token与id是否匹配

	// 匹配, 根据userID查找数据库，按上传时间排序
	videoList, err := repository.UserVideoList(userID)
	if err != nil {
		// 错误处理，返回空列表
		return videoList, err
	}
	// 将视频列表中填充author信息
	for i := range *videoList {

		user_id := (*videoList)[i]["user_id"].(uint64)

		var author *map[string]interface{}
		author, err = repository.AuthorInfo(user_id)
		if err != nil {
			// 错误处理

		}
		(*author)["id"] = (*author)["user_id"]
		(*author)["name"] = (*author)["user_name"]
		(*author)["is_follow"] = false

		(*videoList)[i]["author"] = author

		// 将相对url转为完整url
		(*videoList)[i]["play_url"] = server_ip + (*videoList)[i]["play_url"].(string)
		(*videoList)[i]["cover_url"] = server_ip + (*videoList)[i]["cover_url"].(string)
	}

	return videoList, nil
	// 不匹配
}
