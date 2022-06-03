package service

import (
	"DouYin/repository"
	"errors"
)

// 关注列表的用户信息
type FollowResponse struct {
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// 关注操作
func RelationAction(userID uint64, token string, toUserID uint64, action bool) error {

	// 验证token
	tokenID, err := Token2ID(token)
	if err != nil {
		return err
	}

	if userID != tokenID {
		return errors.New("token.user_id 与 user_id不同")
	}

	err = repository.Action(userID, toUserID, action)
	return err
}

// 获取关注列表
func FollowList(userID uint64) (*[]FollowResponse, error) {
	var followList []FollowResponse
	followUserIDList, err := repository.FollowList(userID)
	if err != nil {
		return nil, err
	}

	for i := range *followUserIDList {
		followUserID := (*followUserIDList)[i].ToUserID

		// 获取其他信息
		userInfo, err := repository.UserInfo(int64(userID))
		if err != nil {
			return nil, err
		}

		followList = append(followList, FollowResponse{
			ID:            followUserID,
			Name:          userInfo.UserName,
			FollowCount:   userInfo.FollowCount,
			FollowerCount: userInfo.FollowerCount,
			IsFollow:      true,
		})
	}

	return &followList, nil
}

// 获取粉丝列表
func FollowerList(userID uint64) (*[]FollowResponse, error) {
	var followerList []FollowResponse
	followerUserIDList, err := repository.FollowerList(userID)
	if err != nil {
		return nil, err
	}

	for i := range *followerUserIDList {
		followUserID := (*followerUserIDList)[i].ToUserID

		// 获取其他信息
		userInfo, err := repository.UserInfo(int64(userID))
		if err != nil {
			return nil, err
		}

		isFollow, err := repository.IsFollower(userID, followUserID)
		if err != nil {
			return nil, err
		}

		followerList = append(followerList, FollowResponse{
			ID:            followUserID,
			Name:          userInfo.UserName,
			FollowCount:   userInfo.FollowCount,
			FollowerCount: userInfo.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	return &followerList, nil
}
