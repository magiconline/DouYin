package repository

// "DouYin/json"
import (
	"fmt"
	"strconv"
)

type User struct {
	// gorm.Model
	UserId        uint64 `gorm:"column:user_id; primary_key; AUTO_INCREMENT"`
	UserName      string `gorm:"column:user_name"`
	Password      string `gorm:"column:password"`
	FollowCount   uint64 `gorm:"column:follow_count"`
	FollowerCount uint64 `gorm:"column:follower_count"`
}

//根据email pwd查找用户
// func FindUserbyEmailandPwd(email, pwd string) (*User, error) {
// 	var user User
// 	err := DB.Table("user").Where("email = ? and password = ?", email, pwd).First(&user).Error
// 	return &user, err
// }

func FindUserbyNameandPwd(name, pwd string) (*User, int, error) {
	var user User
	res := DB.Table("user").Where("user_name = ? and password = ?", name, pwd).Limit(1).Find(&user)
	return &user, int(res.RowsAffected), res.Error
}

//根据email查找用户
// func FindUserbyEmail(email string) (*User, error) {
// 	var user User
// 	err := DB.Table("user").Where("email = ? ", email).First(&user).Error
// 	return &user, err
// }

func FindUserbyName(name string) (int, error) {
	var user User
	res := DB.Table("user").Where("user_name = ? ", name).Limit(1).Find(&user)
	return int(res.RowsAffected), res.Error
}

//根据userid查找用户
func FindUserbyID(userid string) (int, error) {
	var user User
	id, _ := strconv.Atoi(userid)
	res := DB.Table("user").Where("user_id = ?", id).Limit(1).Find(&user)
	return int(res.RowsAffected), res.Error
}

// 创建用户
func CreateUser(username, pwd string) (*User, error) {
	// _, err1 := FindUserbyEmail(email)
	// if err1 == nil {
	user := User{
		UserName: username,
		Password: pwd,
	}
	err := DB.Table("user").Create(&user).Error
	return &user, err
}

// 为关注/粉丝列表查询用户信息(user_name, follow_count, follower_count)
// 使用redis缓存
func UserInfo(userID uint64) (*User, error) {
	var result User

	// 查询缓存
	key := fmt.Sprintf("user_%v", userID)
	rdbResult, err := RDB.HGetAll(CTX, key).Result()
	if err == nil && len(rdbResult) != 0 {
		// 缓存找到
		result.UserName = rdbResult["user_name"]
		follow_count, _ := strconv.ParseUint(rdbResult["follow_count"], 10, 64)
		result.FollowCount = follow_count
		follower_count, _ := strconv.ParseUint(rdbResult["follower_count"], 10, 64)
		result.FollowerCount = follower_count
		return &result, nil
	}

	// 没有找到缓存
	err = DB.Table("user").Where(User{UserId: userID}).Select("user_name", "follow_count", "follower_count").Limit(1).Find(&result).Error

	// 更新缓存
	RDB.HSet(CTX, key, "user_name", result.UserName, "follow_count", result.FollowCount, "follower_count", result.FollowerCount)

	return &result, err
}
