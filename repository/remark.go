package repository

import (
	"strconv"
	"sync"
	"time"
)

type Remark struct {
	CommentId 		int64	`gorm:"column:comment_id"`		//评论id
	VideoId         int64	`gorm:"column:video_id"`		//视频id
	UserId			int64	`gorm:"column:user_id"`		//发出该评论的用户id
	ActionType      int32	`gorm:"column:action_type"`		//1-发布评论，2-删除评论
	CommentText		string	`gorm:"column:comment_text"`		//用户填写的评论内容
	CreateTime		time.Time	`gorm:"column:create_time"`		//评论时间
}

type Remark1 struct {
	CommentId 		int64	`gorm:"column:comment_id"`		//评论id
	User			User	`gorm:"column:user_id"`		//发出该评论的用户id
	Content			string	`gorm:"column:comment_text"`		//用户填写的评论内容
	CreateDate		time.Time	`gorm:"column:create_time"`		//评论时间
}


//设置表名，可以通过给struct类型定义 TableName函数，返回当前struct绑定的mysql表名是什么
func (u Remark) TableName() string {
	//绑定MYSQL表名为users
	return "remark"
}

type RemarkDao struct {
}

var remarkDao *RemarkDao
var remarkOnce sync.Once

func NewRemarkDaoINstance()*RemarkDao  {
	remarkOnce.Do(
		func() {
			remarkDao = &RemarkDao{}
		})
	return remarkDao

}

//根据视频id查询评论
func  (*RemarkDao)QueryByVideoId(videoId uint64) ([]map[string]interface{},error) {
	var  remarkList  []map[string]interface{}
	result := DB.Table("remark").Where("video_id = ?",videoId).Order("create_time desc").Find(&remarkList).Error
	return remarkList,result
}

//根据评论id查询评论
func  (*RemarkDao)QueryByCommentId(commentId uint64) ([]map[string]interface{},error) {
	var  remarkList  []map[string]interface{}
	result := DB.Table("remark").Where("comment_id = ?",commentId).Order("create_time desc").Find(&remarkList)

	return remarkList,result.Error
}


//插入评论,计算videoid 下的评论总条数,同时更新到video列表
func (*RemarkDao)InsertByCommentIDAndVideo(remarkdata *Remark,videoId uint64) (error) {
	var count int64
	//开启事务
	tx := DB.Begin()
	err := tx.Create(remarkdata).Error;
	err = tx.Table("video").Where("video_id = ?", videoId).Update("comment_count", count+1).Error
	if err != nil {
		tx.Rollback()
		return nil
	}
	return nil

}


// 根据videoID查找所有评论
func VideoCommentList(videoId uint64) (*[]map[string]interface{}, error) {
	var commentList []map[string]interface{}
	err := DB.Table("remark").Where("video_id = ?",videoId).Order("create_time desc").Find(&commentList).Error

	return &commentList, err
}


// 根据videoid和userid和comment_text删除评论,comment_id为主键
func DeleteByVdUdAndContent(videoId uint64,userId uint64,content string)(error){
	var count int64
	tx := DB.Begin()
	err :=tx.Table("remark").Where("video_id = ?",videoId).Where("user_id = ?", userId).Where("comment_text = ?",content).Delete(&Remark{})
	err = tx.Table("video").Where("video_id = ?", videoId).Update("comment_count", count+1)
	if err !=nil {
		tx.Rollback()
		return nil
	}
	return err.Error
}


//
func CountCommentlist(videoId uint64)(err error) {
	var count int64
	err = DB.Table("video").Where("video_id = ?", videoId).Update("comment_count", count+1).Error
	if err !=nil {
		return err
	}
	return err
}

//根据userid查找用户
func FindUserbyIDComment(userid string) (*User, error) {
	var user User
	id, _ := strconv.Atoi(userid)
	err := DB.Table("user").Where("user_id = ?", id+1).First(&user).Error
	return &user, err
}

