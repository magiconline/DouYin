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
	ID         uint64       `json:"id"`
	User       UserResponse `json:"user"`
	Content    string       `json:"content"`
	CreateDate time.Time    `json:"create_date"`
}

type CommentActionResponse CommentResponse

//获取userid对应的相应值
func UserInfo(userId uint64) (*UserResponse, error) {
	user, err := repository.UserInfo(int64(userId))
	if err != nil {
		return nil, err
	}
	return &UserResponse{ID: uint64(user.UserId), Name: user.UserName, FollowCount: 0, FollowerCount: 0, IsFollow: false}, nil
}

//将评论内容插入数据库
func InsertByCommentIdAndVideo(remark repository.Remark, videoId uint64) (int64, error) {
	cid, err := repository.NewRemarkDaoINstance().InsertByCommentIDAndVideo(&remark, videoId)
	if err != nil {
		return cid, err
	}
	return cid, err
}

// func QueryByCommentId(commentId int64) ([]map[string]interface{}, error) {

//获得新插入评论的响应值
func NewInsetRemark(token string, videoId uint64, remark repository.Remark) (*CommentActionResponse, error) {
	//检查token
	if token != "" {
		_, err := Token2ID(token)
		if err != nil {
			return nil, err
		}
	}
	user_id, _ := Token2ID(token)
	user, err := UserInfo(user_id)
	if err != nil {
		return nil, err
	}

	cid, err := InsertByCommentIdAndVideo(remark, videoId)
	if err != nil {
		return nil, err
	}

	var response CommentActionResponse
	//将评论列表填充user信息
	response = CommentActionResponse{
		ID:         uint64(cid),
		User:       *user,
		Content:    remark.CommentText,
		CreateDate: remark.CreateTime,
	}

	return &response, nil

}

//获得videoId的所有评论列表
func VideoList(token string, videoid uint64) (*[]CommentActionResponse, error) {
	//检查token
	if token != "" {
		_, err := Token2ID(token)
		if err != nil {
			return nil, err
		}
	}
	user_id, _ := Token2ID(token)
	user, err := UserInfo(user_id)
	if err != nil {
		return nil, err
	}

	//根据videoId查找数据库中对应所有评论
	commentList, err := repository.VideoCommentList(videoid)
	if err != nil {
		return nil, err
	}

	var response []CommentActionResponse
	//将评论列表填充user信息
	for i := range *commentList {
		response_i := CommentActionResponse{
			ID:         uint64((*commentList)[i]["comment_id"].(int32)),
			User:       *user,
			Content:    (*commentList)[i]["comment_text"].(string),
			CreateDate: (*commentList)[i]["create_time"].(time.Time),
		}
		response = append(response, response_i)
	}
	return &response, nil

}

//func DeleteContent(VideoID uint64,UserId uint64,Content string) (error) {
//	err:=repository.DeleteByVdUdAndContent(VideoID,UserId,Content)
//
//	return err
//}
//根据视频id和commentId删除评论
func DeleteByCommentID(commentId uint64, videoId uint64) error {
	err := repository.DeleteByComentID(commentId, videoId)
	return err
}
