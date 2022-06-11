package controller

import (
	"DouYin/logger"
	"DouYin/repository"
	"DouYin/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//登录
func UserLogin(c *gin.Context) *gin.H {
	// return nil
	//1.收到请求参数，取出参数中的email
	username := c.Query("username")
	pwd := c.Query("password")

	//2.对请求参数进行表单验证，以保证数据库的安全---待完成

	// 3.数据库查询用户名
	_, err1 := repository.FindUserbyName(username)
	if err1 == gorm.ErrRecordNotFound {
		logger.Println(err1.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  "用户名不存在,请注册",
		}
	} else if err1 == nil {
		user, err := repository.FindUserbyNameandPwd(username, pwd)
		//如果数据库查不到或者密码错误
		if err == gorm.ErrRecordNotFound {
			logger.Println(err.Error())
			return &gin.H{
				"status_code": 1,
				"status_msg":  "用户密码错误",
			}
		} else if err != nil {
			logger.Println(err.Error())
			return &gin.H{
				"status_code": 1,
				"status_msg":  "数据库查询出错，请重新登录",
			}
		}

		//token
		token, err := service.GenerateToken(uint(user.UserId))
		if err != nil {
			logger.Println(err.Error())
			return &gin.H{
				"status_code": 1,
				"status_msg":  "token生成失败请重新登录",
			}
		}

		//查到了就返回json
		return &gin.H{
			"status_code": 0,
			"status_msg":  "SUCCESS",
			"user_id":     user.UserId,
			"token":       token,
		}
	}
	logger.Println(err1.Error())
	return &gin.H{
		"status_code": 1,
		"status_msg":  "数据库查询出错",
	}
}

//注册
func UserRegister(c *gin.Context) *gin.H {
	username := c.Query("username")
	pwd := c.Query("password")
	_, err1 := repository.FindUserbyName(username)
	//若数据库中不存在此email
	// if err == gorm.ErrRecordNotFound {
	if err1 == gorm.ErrRecordNotFound {
		u, err := repository.CreateUser(username, pwd)
		//若创建新用户不成功
		if err != nil {
			return &gin.H{
				"status_code": 1,
				"status_msg":  "数据库创建用户操作错误",
			}
		}
		token, err := service.GenerateToken(uint(u.UserId))
		if err != nil {
			return &gin.H{
				"status_code": 1,
				"status_msg":  "token生成失败请重新登录",
			}
		}
		return &gin.H{
			"status_code": 0,
			"status_msg":  "success",
			"user_id":     u.UserId,
			"token":       token,
		}
	} else {
		//若查找返回的err不是不存在而是别的
		if err1 != nil {
			return &gin.H{
				"status_code": 1,
				"status_msg":  "数据库查询用户操作错误",
			}
		}
		if err1 == nil {
			return &gin.H{
				"status_code": 1,
				"status_msg":  "该用户名已被注册",
			}
		}
	}
	return &gin.H{
		"status_code": 1,
		"status_msg":  "注册失败",
	}

}

func UserInfo(c *gin.Context) *gin.H {
	userid := c.Query("user_id")
	token := c.Query("token")

	type UserInformation struct {
		UserId        int64  `gorm:"column:user_id"`
		UserName      string `gorm:"column:user_name"`
		FollowCount   int    `gorm:"column:follow_count"`
		FollowerCount int    `gorm:"column:follower_count"`
	}

	user, err := repository.FindUserbyID(userid)
	if err == gorm.ErrRecordNotFound {
		return &gin.H{
			"status_code": 1,
			"status_msg":  "用户不存在",
		}
	}

	// if token == "" {
	// 	return &gin.H{
	// 		"status_code": 1,
	// 		"status_msg":  "token不存在",
	// 	}
	// }

	_, err = service.ParseToken(token)
	if err != nil {
		return &gin.H{
			"status_code": 1,
			"status_msg":  "token解析失败",
		}
	}

	return &gin.H{
		"status_code": 0,
		"status_msg":  "sucess",
		"user": gin.H{
			"id":             user.UserId,
			"name":           user.UserName,
			"follow_count":   user.FollowCount,
			"follower_count": user.FollowerCount,
		},
	}
}
