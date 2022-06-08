package service

import (
	"DouYin/repository"
	"errors"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
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
func RelationAction(token string, toUserID uint64, action bool) error {
	var err error

	// 验证token
	userID, err := Token2ID(token)
	if err != nil {
		return err
	}

	// 生成锁信息
	k := fmt.Sprintf("%d_%d", userID, toUserID)
	v := uuid.NewString()

	// 获得锁
	if err = repository.GetRedisLock(k, v, 10*time.Second); err != nil {
		return err
	}

	// 释放锁
	defer repository.DeleteRedisLock(k, v)

	// 开启事务更新relation表和user表
	tx := repository.DB.Begin()
	// 修改relation
	err = repository.Action(tx, userID, toUserID, action)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新user表关注数和粉丝数
	var value int
	if action {
		value = 1
	} else {
		value = -1
	}

	err1 := repository.ChangeFollowCount(tx, userID, value)
	err2 := repository.ChangeFollowerCount(tx, toUserID, value)
	if err1 != nil || err2 != nil {
		tx.Rollback()
		return errors.New("user.follow_count | user.follower_count 更新失败")
	}

	tx.Commit()
	return nil
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
		userInfo, err := repository.UserInfo(int64(followUserID))
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
		followerUserID := (*followerUserIDList)[i].ToUserID

		// 获取其他信息
		userInfo, err := repository.UserInfo(int64(followerUserID))
		if err != nil {
			return nil, err
		}

		isFollow, err := repository.IsFollower(userID, followerUserID)
		if err != nil {
			return nil, err
		}

		followerList = append(followerList, FollowResponse{
			ID:            followerUserID,
			Name:          userInfo.UserName,
			FollowCount:   userInfo.FollowCount,
			FollowerCount: userInfo.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	return &followerList, nil
}
