package controller

import (
	"DouYin/logger"
	"DouYin/repository"
	"DouYin/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

//登录
func UserLogin(c *gin.Context) *gin.H {
	// return nil
	//1.收到请求参数，取出参数中的email
	username := c.Query("username")
	pwd := c.Query("password")

	//2.对请求参数进行表单验证，以保证数据库的安全---待完成

	// 3.数据库查询用户名
	numOfRowsAffected, err1 := repository.FindUserbyName(username)
	if numOfRowsAffected == 0 {
		return &gin.H{
			"status_code": 1,
			"status_msg":  "用户名不存在,请注册",
		}
	} else if err1 == nil {
		user, num1OfRowsAffected, err := repository.FindUserbyNameandPwd(username, pwd)
		//如果数据库查不到或者密码错误
		if num1OfRowsAffected == 0 {
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
		token, err := service.GenerateToken(user.UserId)
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
	numOfRowsAffected, err1 := repository.FindUserbyName(username)
	//若数据库中不存在此email
	// if err == gorm.ErrRecordNotFound {
	if numOfRowsAffected == 0 {
		u, err := repository.CreateUser(username, pwd)
		//若创建新用户不成功
		if err != nil {
			return &gin.H{
				"status_code": 1,
				"status_msg":  "数据库创建用户操作错误",
			}
		}
		token, err := service.GenerateToken(u.UserId)
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
		if numOfRowsAffected > 0 {
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
	useridStr := c.Query("user_id")
	token := c.Query("token")
	UserID, err := strconv.ParseUint(useridStr, 10, 64)
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		}
	}

	_, err1 := service.ParseToken(token)
	if err1 != nil {
		logger.Println(err1.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err1.Error(),
		}
	}

	user, err := repository.UserInfo(UserID)
	if err != nil {
		logger.Println(err.Error())
		return &gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
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
