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
	FollowCount   uint64 `json:"follow_count"`
	FollowerCount uint64 `json:"follower_count"`
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

	// 如果自己关注自己则返回错误
	if userID == toUserID {
		return errors.New("禁止关注自己")
	}

	// 检查userID是否存在
	exist, err := repository.UserIDExist(userID, toUserID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("用户user_id不存在")
	}

	// 生成锁信息
	k := fmt.Sprintf("relation_%d_%d", userID, toUserID)
	v := uuid.NewString()

	// 获得锁
	if err = repository.GetRedisLock(k, v, 10*time.Second); err != nil {
		return err
	}

	// 释放锁
	defer repository.DeleteRedisLock(k, v)

	// 检测重复关注/取消关注
	if isFollow, _ := repository.IsFollower(userID, toUserID); isFollow == action {
		if isFollow {
			return errors.New("已关注，禁止重复操作")
		} else {
			return errors.New("已取消关注，禁止重复操作")
		}
	}

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

	// 删除redis缓存
	key1 := fmt.Sprintf("user_%v", userID)
	key2 := fmt.Sprintf("user_%v", toUserID)
	_, err = repository.RDB.Del(repository.CTX, key1, key2).Result()
	if err != nil {
		return err
	}

	return nil
}

// 获取userID的关注列表
// curUserID用来获取IsFollow信息
func FollowList(curUserID uint64, userID uint64) (*[]FollowResponse, error) {
	var followList []FollowResponse
	followUserIDList, err := repository.FollowList(userID)
	if err != nil {
		return nil, err
	}

	for i := range *followUserIDList {
		followUserID := (*followUserIDList)[i].ToUserID

		// 获取其他信息
		userInfo, err := repository.UserInfo(followUserID)
		if err != nil {
			return nil, err
		}

		isFollow, err := repository.IsFollower(curUserID, userID)
		if err != nil {
			return nil, err
		}

		followList = append(followList, FollowResponse{
			ID:            followUserID,
			Name:          userInfo.UserName,
			FollowCount:   userInfo.FollowCount,
			FollowerCount: userInfo.FollowerCount,
			IsFollow:      isFollow,
		})
	}

	return &followList, nil
}

// 获取userID的粉丝列表
// curUserID 用来查询curUserID是否已关注userID
func FollowerList(curUserID uint64, userID uint64) (*[]FollowResponse, error) {
	var followerList []FollowResponse
	followerUserIDList, err := repository.FollowerList(userID)
	if err != nil {
		return nil, err
	}

	for i := range *followerUserIDList {
		followerUserID := (*followerUserIDList)[i].UserID

		// 获取其他信息
		userInfo, err := repository.UserInfo(followerUserID)
		if err != nil {
			return nil, err
		}

		isFollow, err := repository.IsFollower(curUserID, followerUserID)
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
