package repository

// "DouYin/json"
import (
	"strconv"
)

type User struct {
	// gorm.Model
	UserId        int64  `gorm:"column:user_id; primary_key; AUTO_INCREMENT"`
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

func FindUserbyNameandPwd(name, pwd string) (*User, error) {
	var user User
	err := DB.Table("user").Where("user_name = ? and password = ?", name, pwd).First(&user).Error
	return &user, err
}

//根据email查找用户
// func FindUserbyEmail(email string) (*User, error) {
// 	var user User
// 	err := DB.Table("user").Where("email = ? ", email).First(&user).Error
// 	return &user, err
// }

func FindUserbyName(name string) (*User, error) {
	var user User
	err := DB.Table("user").Where("user_name = ? ", name).First(&user).Error
	return &user, err
}

//根据userid查找用户
func FindUserbyID(userid string) (*User, error) {
	var user User
	// fmt.Println("string = ", userid)
	id, _ := strconv.Atoi(userid)
	// fmt.Println("int = ", id)
	err := DB.Table("user").Where("user_id = ?", id).First(&user).Error
	return &user, err
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
func UserInfo(userID int64) (*User, error) {
	var result User

	err := DB.Table("user").Where(User{UserId: userID}).Select("user_name", "follow_count", "follower_count").Take(&result).Error

	return &result, err
}
