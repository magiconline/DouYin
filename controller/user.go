package controller

import (
	"github.com/gin-gonic/gin"
)

//登录注册
func UserStore(c *gin.Context) *interface{} {
	return nil
	//1.收到请求参数，取出参数中的username
	//2.对请求参数进行验证，有错误就返回json
	//3.数据库查询用户username

}

//根据用户名查找用户
// func FindUserByName(name string) (*User, error) {
// 	var user User
// 	// err :=
// 	// return &user, err
// 	return nil
// }

//创建用户
// func CreateUser(name string) (*User, error) {

// }
