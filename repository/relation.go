package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Relation struct {
	ID       uint64 `gorm:"primaryKey"`
	UserID   uint64 `gorm:"notNULL;index"` // 用户
	ToUserID uint64 `gorm:"notNULL;index"` // 被关注的用户
	Relation bool   `gorm:"notNULL"`
}

// 获取redis锁
// 如果锁已被占用则休眠0.1s后继续查询
func GetRedisLock(k string, v string, expire time.Duration) error {
	val := false
	var err error

	for !val {
		val, err = RDB.SetNX(CTX, k, v, expire).Result()
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// 释放redis锁
// 可重复释放
func DeleteRedisLock(k string, v string) {
	result, err := RDB.Get(CTX, k).Result()
	if err != nil {
		fmt.Printf("redis err: %s, %s 释放失败\n", err.Error(), k)
	}

	if result == v {
		err = RDB.Del(CTX, k).Err()
		if err != nil {
			fmt.Printf("redis err: %s, %s 释放失败\n", err.Error(), k)
		}
	}
}

// 关注操作
// action=true表示关注操作
// action=false表示取消关注操作
func Action(tx *gorm.DB, userID uint64, toUserID uint64, action bool) error {
	relation := Relation{UserID: userID, ToUserID: toUserID}
	var err error
	if action {
		// 关注操作

		// 查询是否已创建relation记录
		result := tx.Where(&relation).Limit(1).Find(&Relation{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			// 未创建relation，创建新relation记录
			relation.Relation = true
			err = tx.Create(&relation).Error
			if err != nil {
				return nil
			}
		} else {
			// 已创建relation，修改relation记录
			err = tx.Model(&relation).Where("user_id = ? AND to_user_id = ?", userID, toUserID).Update("relation", true).Error
		}

		return err
	} else {
		// 取消关注操作
		err := tx.Model(&relation).Where("user_id = ? AND to_user_id = ?", userID, toUserID).Update("relation", false).Error
		return err
	}
}

// 更新user表follow_count
func ChangeFollowCount(tx *gorm.DB, userID uint64, value int) error {
	err := tx.Table("user").Where("user_id = ?", userID).Update("follow_count", gorm.Expr("follow_count + (?)", value)).Error
	return err
}

// 更新user表的follower_count
func ChangeFollowerCount(tx *gorm.DB, userID uint64, value int) error {
	err := tx.Table("user").Where("user_id = ? ", userID).Update("follower_count", gorm.Expr("follower_count + (?)", value)).Error
	return err
}

// 关注列表
func FollowList(userID uint64) (*[]Relation, error) {
	var results []Relation
	err := DB.Where(&Relation{UserID: userID, Relation: true}).Select("to_user_id").Find(&results).Error

	return &results, err
}

// 粉丝列表
func FollowerList(userID uint64) (*[]Relation, error) {
	var results []Relation
	err := DB.Where(&Relation{ToUserID: userID, Relation: true}).Select("user_id").Find(&results).Error

	return &results, err
}

// 判断userID是否关注了toUserID
func IsFollower(userID uint64, toUserID uint64) (bool, error) {
	relation := &Relation{}
	err := DB.Where(&Relation{UserID: userID, ToUserID: toUserID, Relation: true}).Limit(1).Find(&relation).Error
	return relation.Relation, err
}

func UserIDExist(userID uint64, toUserID uint64) (bool, error) {
	var count int64 = 0
	err := DB.Table("user").Where("user_id  IN ?", []uint64{userID, toUserID}).Count(&count).Error
	return count == 2, err
}
