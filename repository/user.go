package repository

import (
	"DouYin/json"

	// "github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserId        int64  `gorm:"column:user_id"`
	Email         string `gorm:"column:user_email"`
	UserName      string `gorm:"column:user_name"`
	Password      string `gorm:"column:password"`
	Token         string `gorm:"column:token"`
	FollowCount   int    `gorm:"column:follow_count"`
	FollowerCount int    `gorm:"column:follower_count"`
}

//根据email查找用户
func FindUserbyEmail(email string) (*User, error) {
	var user User
	err := DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// 创建用户
func CreateUser(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := DB.Create(&user).Error
	return &user, err
}

//登录注册
func UserStore(c *gin.Context) {
	// return nil
	//1.收到请求参数，取出参数中的email
	email := c.PostForm("Email")
	// vCode := c.PostForm("code")

	//2.对请求参数进行表单验证，以保证安全
	// validate := validation.Validation{}
	// validate.Required(email, "Email").Message("邮箱有误")

	//3.数据库查询用户username
	user, err := FindUserbyEmail(email)
	//如果数据库查不到就需要注册
	if err == gorm.ErrRecordNotFound {
		user, err = CreateUser(email)
		if err != nil {
			json.ResponseWithJson(json.ERROR, "数据库操作错误", c)
			return
		}
	} else {
		if err != nil {
			json.ResponseWithJson(json.ERROR, "数据库操作错误", c)
			return
		}
	}
	//查到了就返回json
	json.ResponseWithJson(json.SUCCESS, gin.H{
		"User": user,
	}, c)
}
