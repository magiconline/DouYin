package controller

import (
	"DouYin/logger"
	"DouYin/repository"
	"DouYin/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func strToInt(a string) int {
	t, err := strconv.Atoi(a)
	if err != nil {
		logger.Logger.Println("error:", err)
	}
	return t
}

/*
	发表评论及删除
*/
func Leave_remark(c *gin.Context) *gin.H {
	token := c.Query("token")
	userId, _ := service.Token2ID(token) //非空
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")   //只能为1/2
	commentText := c.Query("comment_text") //为字符串，不需要过多校验
	createTime := time.Now()

	if token == "" {
		return &gin.H{
			"status_code": 1,
			"status_msg":  "用户未登录",
		}
	}

	//类型转换
	videoIdInt := strToInt(videoId)
	actionTypeInt := strToInt(actionType)

	remark := repository.Remark{
		VideoId:     int64(videoIdInt),
		UserId:      int64(userId),
		ActionType:  int32(actionTypeInt),
		CommentText: commentText,
		CreateTime:  createTime,
	}

	var comment1 *service.CommentActionResponse
	//插入评论
	if actionTypeInt == 1 {
		//插入评论不需要检查该字段是否已经评论过
		//插入数据库
		comment, err := service.NewInsetRemark(token, uint64(videoIdInt), remark)
		if err != nil {
			return &gin.H{
				"status_code": 1,
				"status_msg":  err.Error(),
				"comment":     "数据库插入失败",
			}
		}
		comment1 = comment
		if err != nil {
			//错误返回
			return &gin.H{
				"status_code": 1,
				"status_msg":  err.Error(),
			}
		}
	}

	if actionTypeInt == 2 {
		//获取评论参数
		commentId := c.Query("comment_id")
		err := service.DeleteByCommentID(uint64(strToInt(commentId)), uint64(videoIdInt))
		if err != nil {
			return &gin.H{
				"status_code": 1,
				"status_msg":  "数据库删除失败",
			}
		}
		return &gin.H{
			"status_code": 0,
			"status_msg":  "删除成功",
			"comment":     0,
		}
	}

	return &gin.H{
		"status_code": 0,
		"status_msg":  "插入成功",
		"comment":     comment1,
	}
}

/*
	查看视频的所有评论，按发布时间倒序
*/
func View_video_remark(c *gin.Context) *gin.H {
	token := c.Query("token")
	videoId := c.Query("video_id") //非空

	//类型转换
	videoIdInt, err := strconv.Atoi(videoId)
	if err != nil {
		logger.Logger.Println("error:", err)
	}
	/**查询所有评论，放到commentList中*/
	comment_list, err := service.VideoList(token, uint64(videoIdInt))
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}
	//service.GetAllComment(token, uint64(videoIdInt))

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

///*
//	删除评论
// */
//func Delete_video_remark(VideoID uint64,UserId uint64,Content string)(error) {
//	err :=service.DeleteContent(VideoID,UserId,Content)
//	return err
//}
