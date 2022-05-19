package controller

// import (
// 	"DouYin/json"
// 	"DouYin/repository"
// 	"gorm.io/gorm"
// 	"github.com/gin-gonic/gin"
// )
// type User struct {
// 	gorm.Model
// 	UserId        int64  `gorm:"column:user_id"`
// 	Email         string `gorm:"column:user_email"`
// 	UserName      string `gorm:"column:user_name"`
// 	Password      string `gorm:"column:password"`
// 	Token         string `gorm:"column:token"`
// 	FollowCount   int    `gorm:"column:follow_count"`
// 	FollowerCount int    `gorm:"column:follower_count"`
// }
// //登录注册
// func UserStore(c *gin.Context) {
// 	// return nil
// 	//1.收到请求参数，取出参数中的username
// 	// user_name := c.PostForm("UserName")
// 	email := c.PostForm("Email")
// 	vCode := c.PostForm("code")

// 	json.ResponseWithJson(json.SUCCESS, gin.H{
// 		"email": email,
// 		"code":  vCode,
// 	}, c)
// 	//2.对请求参数进行验证，有错误就返回json
// 	//3.数据库查询用户username

// }

// //根据email查找用户
// // func FindUserByName(name string) (*User, error) {
// // 	var user User
// // 	// err :=
// // 	// return &user, err
// // 	return nil
// // }
// func FindUserbyEmail(email string) (*User, error) {
// 	var user User
// 	err := DB.Where("email = ?", email).First(&user).error
// 	return &user, err
// }

// //创建用户
// // func CreateUser(name string) (*User, error) {

// // }
