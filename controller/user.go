package controller

import (
	"DouYin/repository"
	"DouYin/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//登录
func UserLogin(c *gin.Context) *gin.H {
	// return nil
	//1.收到请求参数，取出参数中的email
	email := c.Query("Email")
	pwd := c.Query("Password")

	//2.对请求参数进行表单验证，以保证数据库的安全---待完成

	// 3.数据库查询用户email
	_, err1 := repository.FindUserbyEmail(email)
	if err1 == gorm.ErrRecordNotFound {
		return &gin.H{
			"code": 1,
			"msg":  "用户邮箱不存在,请注册",
		}
	} else if err1 == nil {
		user, err := repository.FindUserbyEmailandPwd(email, pwd)
		//如果数据库查不到或者密码错误
		if err == gorm.ErrRecordNotFound {
			return &gin.H{
				"code": 1,
				"msg":  "用户密码错误",
			}
		} else if err != nil {
			return &gin.H{
				"code": 1,
				"msg":  "数据库查询出错，请重新登录",
			}
		}

		//token
		token, err := service.GenerateToken(uint(user.UserId))
		if err != nil {
			return &gin.H{
				"code": 1,
				"msg":  "token生成失败请重新登录",
			}
		}

		//查到了就返回json
		return &gin.H{
			"code":    0,
			"msg":     "SUCCESS",
			"user_id": user.UserId,
			"token":   token,
		}
	}
	return &gin.H{
		"code": 1,
		"msg":  "数据库查询出错",
	}
}

//注册
func UserRegister(c *gin.Context) *gin.H {
	email := c.Query("Email")
	pwd := c.Query("Password")
	_, err1 := repository.FindUserbyEmail(email)
	//若数据库中不存在此email
	// if err == gorm.ErrRecordNotFound {
	if err1 == gorm.ErrRecordNotFound {
		u, err := repository.CreateUser(email, pwd)
		//若创建新用户不成功
		if err != nil {
			return &gin.H{
				"code": 1,
				"msg":  "数据库创建用户操作错误",
			}
		}
		token, err := service.GenerateToken(uint(u.UserId))
		if err != nil {
			return &gin.H{
				"code": 1,
				"msg":  "token生成失败请重新登录",
			}
		}
		return &gin.H{
			"code":    0,
			"msg":     "success",
			"user_id": u.UserId,
			"token":   token,
		}
	} else {
		//若查找返回的err不是不存在而是别的
		if err1 != nil {
			return &gin.H{
				"code": 1,
				"msg":  "数据库查询用户操作错误",
			}
		}
		if err1 == nil {
			return &gin.H{
				"code": 1,
				"msg":  "该邮箱已被注册",
			}
		}
	}
	return &gin.H{
		"code": 1,
		"msg":  "注册失败",
	}

}
func UserInfo(c *gin.Context) *gin.H {
	userid := c.Query("UserId")
	token := c.Query("token")

	user, err := repository.FindUserbyID(userid)
	if err == gorm.ErrRecordNotFound {
		return &gin.H{
			"code": 1,
			"msg":  "用户不存在",
		}
	}

	if token == "" {
		return &gin.H{
			"code": 1,
			"msg":  "token不存在",
		}
	}

	_, err = service.ParseToken(token)
	if err != nil {
		return &gin.H{
			"code": 1,
			"msg":  "token解析失败",
		}
	}

	return &gin.H{
		"code":           0,
		"msg":            "sucess",
		"id":             user.UserId,
		"name":           user.UserName,
		"follow_count":   user.FollowCount,
		"follower_count": user.FollowerCount,
	}
}
