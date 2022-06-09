package service

import (
	"DouYin/repository"
	"time"
)

type UserResponse struct {
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint64 `json:"follow_count"`
	FollowerCount uint64 `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type CommentResponse struct {
	ID      	  uint64         `json:"id"`
	User        UserResponse   `json:"user"`
	Content  string     	 `json:"content"`
	CreateDate    time.Time          `json:"create_date"`
}

type Comment struct {

}

type CommentActionResponse CommentResponse

func UserInfo(userId uint64)(*UserResponse,error ) {
	user,err := repository.UserInfo(int64(userId))
	if err != nil{
		return nil,err
	}
	return &UserResponse{ID: uint64(user.UserId), Name:user.UserName, FollowCount: 0,FollowerCount: 0,IsFollow: false},nil


}

func QueryByVideoId(videoId int64) ([]map[string]interface{},error) {

	remarkList, err :=repository.NewRemarkDaoINstance().QueryByVideoId(uint64(videoId))
	if err != nil {
		// 错误处理
		return nil,err
	}
	return remarkList,err
}

func InsertByCommentIdAndVideo(remark repository.Remark,videoId uint64) (error) {
	return repository.NewRemarkDaoINstance().InsertByCommentIDAndVideo(&remark,videoId)
}


func QueryByCommentId(commentId int64) ([]map[string]interface{},error) {

	remarkList, err :=repository.NewRemarkDaoINstance().QueryByCommentId(uint64(commentId))
	if err != nil {
		// 错误处理
		return nil,err
	}
	return remarkList,err
}

//获得videoId的所有评论列表
func VideoList(token string,videoid uint64)(*[]CommentActionResponse,error)  {
	//检查token
	if token != ""{
		_,err := Token2ID(token)
		if err != nil {
			return nil, err
		}
	}
	user_id, _ := Token2ID(token)
	user , err := UserInfo(user_id)
	if err != nil {
		return nil, err
	}

	//根据videoId查找数据库中对应所有评论
	commentList,err:=repository.VideoCommentList(videoid)
	if err != nil {
		return nil, err
	}

	var response []CommentActionResponse
	//将评论列表填充user信息
	for i := range *commentList{
		response_i := CommentActionResponse{
			ID: uint64((*commentList)[i]["comment_id"].(int32)),
			User:        *user,
			Content: 	(*commentList)[i]["comment_text"].(string),
			CreateDate:  (*commentList)[i]["create_time"].(time.Time),
		}
		response = append(response,response_i)
	}
	return &response,nil

}

func DeleteContent(VideoID uint64,UserId uint64,Content string) (error) {
	err:=repository.DeleteByVdUdAndContent(VideoID,UserId,Content)

	return err
}

func CountAllComment(videoId uint64)(err error)  {
	 err=repository.CountCommentlist(videoId)
	return err
}