package controller

import (
	"DouYin/logger"
	"DouYin/repository"
	"DouYin/service"
	"os/user"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RemarkData struct {
	StatusCode uint
	StatusMag  string
}

type Comment struct {
	id         uint
	user       user.User
	content    string
	creat_date time.Time
}

func strToInt(a string) int {
	t, err := strconv.Atoi(a)
	if err != nil {
		logger.Logger.Println("error:", err)
	}
	return t
}

/*
	发表评论
*/
func Leave_remark(c *gin.Context) *gin.H {
	token := c.Query("token")
	userId, _ := service.Token2ID(token) //非空
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")   //只能为1/2
	commentText := c.Query("comment_text") //为字符串，不需要过多校验
	createTime := time.Now()

	//http://localhost:8080/douyin/comment/action/?user_id=1&token=a&video_id=1&action_type=1&comment_text=小草&comment_id=8&parent_id=0

	//类型转换
	videoIdInt := strToInt(videoId)
	actionTypeInt := strToInt(actionType)

	if actionTypeInt == 2 {
		err := Delete_video_remark(uint64(videoIdInt), userId, commentText)
		if err != nil {
			logger.Println(err.Error())
			return &gin.H{
				"status_code": 1,
				"status_msg":  err.Error(),
			}
		}
		return &gin.H{
			"status_code": 0,
			"status_msg":  "success",
			"comment":     "null",
		}
	}

	user1, err := repository.FindUserbyIDComment(string(userId))
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	remark := repository.Remark{
		VideoId:     int64(videoIdInt),
		UserId:      int64(userId),
		ActionType:  int32(actionTypeInt),
		CommentText: commentText,
		CreateTime:  createTime,
	}

	remark1 := repository.Remark1{
		User:       *user1,
		Content:    commentText,
		CreateDate: createTime,
	}
	err = service.InsertByCommentIdAndVideo(remark, uint64(videoIdInt))
	if err != nil {
		//错误返回
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	return &gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"comment":     remark1,
	}

	//1.收到请求参数，获取参数中关于评论的各项参数
	//2.对请求参数进行验证，有错误就返回错误json
	//3.将评论根据video_id插入数据库
}

/*
	查看视频的所有评论，按发布时间倒序
*/
func View_video_remark(c *gin.Context) *gin.H {
	token := c.Query("token")
	videoId := c.Query("video_id") //非空
	//http://localhost:8080/douyin/comment/list/?token=a&video_id=1
	//类型转换
	videoIdInt, err := strconv.Atoi(videoId)
	if err != nil {
		logger.Logger.Println(err.Error())
	}
	/**查询所有评论，放到commentList中*/
	comment_list, err := service.VideoList(token, uint64(videoIdInt))
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	//正常返回
	return &gin.H{
		"status_code":  0,
		"status_msg":   "success",
		"comment_list": comment_list,
	}
	return nil
	//1.收到请求参数，获取参数中关于评论的各项参数
	//2.对请求参数进行验证，有错误就返回错误json
	//3.从数据库中根据video_id查询所有评论条数放到list
	//4.根据list中的时间倒序输出
}

/*
	删除评论
*/
func Delete_video_remark(VideoID uint64, UserId uint64, Content string) error {
	err := service.DeleteContent(VideoID, UserId, Content)

	return err
}
