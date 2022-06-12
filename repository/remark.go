package repository

import (
	"DouYin/logger"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Remark struct {
	Id 				int64		`gorm:"column:comment_id"`			//评论id
	VideoId         int64		`gorm:"column:video_id"`			//视频id
	UserId			int64		`gorm:"column:user_id"`				//发出该评论的用户id
	ActionType      int32		`gorm:"column:action_type"`			//1-发布评论，2-删除评论
	CommentText		string		`gorm:"column:comment_text"`		//用户填写的评论内容
	CreateTime		time.Time	`gorm:"column:create_time"`			//评论时间
}

type Remark1 struct {
	Id 				int64		`gorm:"column:comment_id"`			//评论id
	User			User		`gorm:"column:user_id"`				//发出该评论的用户id
	Content			string		`gorm:"column:comment_text"`		//用户填写的评论内容
	CreateDate		time.Time	`gorm:"column:create_time"`			//评论时间
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

////根据评论id查询评论
//func  (*RemarkDao)QueryByCommentId(commentId uint64) ([]map[string]interface{},error) {
//	var  remarkList  []map[string]interface{}
//	result := DB.Table("remark").Where("comment_id = ?",commentId).Order("create_time desc").Find(&remarkList)
//	return remarkList,result.Error
//}


//插入评论,计算videoid 下的评论总条数,同时更新到video列表
func (*RemarkDao)InsertByCommentIDAndVideo(remarkdata *Remark,videoId uint64) (int64,error) {
	//cid为获取新插入的评论commetn_id
	var cid int64
	cid = 0 //初始化
	var count int64
	count = 0
	//开启事务
	DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(remarkdata).Error ; err != nil{
			logger.Logger.Println("err",err)
			return err
		}
		if result := tx.Table("remark").Where("video_id = ?", videoId).Count(&count) ; result.Error != nil{
			logger.Logger.Println("err",result.Error)
			return result.Error
		}
		if err := tx.Table("video").Where("video_id = ?", videoId).Update("comment_count", count).Error ; err != nil{
			logger.Logger.Println("err",err)
			return err
		}
		//if err := tx.Table("video").Where("video_id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error ; err != nil{
		//	logger.Logger.Println("err",err)
		//	return err
		//}
			return nil
		})
	if err := DB.Table("remark").Select("comment_id").Where("video_id = ?",remarkdata.VideoId).Where("comment_text = ?",remarkdata.CommentText).Where("user_id=?",remarkdata.UserId).Limit(1).Take(&cid).Error; err != nil{
		logger.Logger.Println("err",err)
		return cid,err
	}
	return  cid,nil
}


// 根据videoID查找所有评论
func VideoCommentList(videoId uint64) (*[]map[string]interface{}, error) {
	var commentList []map[string]interface{}
	err := DB.Table("remark").Where("video_id = ?",videoId).Order("create_time desc").Find(&commentList).Error

	return &commentList, err
}


//// 根据videoid和userid和comment_text删除评论,comment_id为主键
//func DeleteByVdUdAndContent(videoId uint64,userId uint64,content string)(error){
//	var count int64
//	tx := DB.Begin()
//	tx.Table("video").Select("comment_count").Where("video_id = ?",videoId).Find(&count)
//	err :=tx.Table("remark").Where("video_id = ?",videoId).Where("user_id = ?", userId).Where("comment_text = ?",content).Unscoped().Delete(&Remark{})
//	err = tx.Table("video").Where("video_id = ?", videoId).Update("comment_count", count-1)
//	if err !=nil {
//		tx.Rollback()
//		return err.Error
//	}
//	//tx.Commit()
//	return err.Error
//}

// 根据commentId 删除评论
func DeleteByComentID(commentId uint64 ,videoId uint64)error{
	var reamrk Remark
	var count int64
	count = 0
	if result := DB.Table("remark").Where("comment_id = ?",commentId).Limit(1).Take(&reamrk); result.Error != nil{
		logger.Logger.Println("err",result.Error)
		return result.Error
	}
	DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("remark").Where("comment_id = ?",commentId).Delete(&Remark{}).Error ; err != nil{
			logger.Logger.Println("err",err)
			return err
		}
		if result := tx.Table("remark").Where("video_id = ?", videoId).Count(&count) ; result.Error != nil{
			logger.Logger.Println("err",result.Error)
			return result.Error
		}
		if err := tx.Table("video").Where("video_id = ?", videoId).Update("comment_count", count).Error ; err != nil{
			logger.Logger.Println("err",err)
			return err
		}
		//if err := tx.Table("video").Where("video_id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error ; err != nil{
		//	logger.Logger.Println("err",err)
		//	return err
		//}
		return nil
	})
	return nil
}



