package main

import (
	"DouYin/logger"
	"DouYin/repository"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var r *gin.Engine

func TestMain(m *testing.M) {
	// 创建日志文件夹
	err := os.Mkdir("./log", 0750)
	if err != nil && !os.IsExist(err) {
		fmt.Println("日志文件夹创建失败")
		os.Exit(1)
	}

	// 初始化项目日志
	err = logger.Init("./log/log.log")
	if err != nil {
		fmt.Println("日志初始化失败,", err)
		os.Exit(1)
	}
	logger.Println("项目日志初始化成功")

	// 初始化静态资源 URL 前缀
	ip, err := getOutBoundIP()
	if err != nil {
		logger.Fatalln("获取本机ip失败,", err.Error())
	}

	port := ":8080"

	setIP(ip, port)

	// 初始化数据库连接
	err = repository.Init()
	if err != nil {
		logger.Fatalln("数据库初始化失败,", err.Error())
	}

	// 设置release模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化gin日志
	file, err := os.OpenFile("./log/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatalln("gin日志初始化失败,", err.Error())
	}
	gin.DefaultWriter = file

	// 创建服务器
	r = gin.Default()

	// 设置路由
	setupRouter(r)

	// 托管静态资源
	r.Static("/static", "./static")

	// 启动服务器
	logger.Println("启动服务器")
	fmt.Printf("启动服务器: %v%v\n", ip, port)

	code := m.Run()

	os.Exit(code)
}

// 测试没有token
func TestFeedWithoutToken(t *testing.T) {
	timeStamp := time.Now().UnixMilli()
	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v", timeStamp))
	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	if int(body["status_code"].(float64)) != 0 {
		t.Errorf("status_code: %v != 0, status_msg: %v", body["status_code"], body["status_msg"])
		t.FailNow()
	}

}

// 测试错误token
func TestFeedWithWrongToken(t *testing.T) {
	timeStamp := time.Now().UnixMilli()
	token := "123"

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v&token=%v", timeStamp, token))

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"].(string), "token contains an invalid number of segments")
}

// 测试过期token
func TestFeedWithExpiredToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiZXhwIjoxNjU0NzAwNzcyLCJpc3MiOiJ6amN5In0.X9VuPerdOP8TNFxVpWY3vLVFPHdVE72un8TiimFMFPk"
	timeStamp := time.Now().UnixMilli()

	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v&token=%v", timeStamp, token))

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)

	result, err := regexp.Match("token is expired by.*", []byte(body["status_msg"].(string)))
	if err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
	if result != true {
		t.Fatalf("status_msg: %v != token is expired by*", body["status_msg"].(string))
		t.FailNow()
	}
}

// 测试访问完所有视频
func TestFeedWithEndTime(t *testing.T) {
	timeStamp := 0
	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v", timeStamp))

	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)

	l := len(body["video_list"].([]interface{}))
	assert.NotEqual(t, l, 0)
	assert.NotEqual(t, body["next_time"], timeStamp)
}

// 正常注册
func TestRegister(t *testing.T) {
	username := "testFirst"
	password := "123456"

	// 检查是否被注册
	var count int64
	err := repository.DB.Table("user").Where("user_name = ?", username).Count(&count).Error
	assert.Equal(t, err, nil)

	// 已注册则删除
	if count != 0 {
		result := repository.DB.Table("user").Where("user_name = ?", username).Delete(&repository.User{})
		assert.Equal(t, result.Error, nil)
		assert.Equal(t, result.RowsAffected, int64(1))
	}

	// 注册
	response, err := http.Post(fmt.Sprintf("http://127.0.0.1:8080/douyin/user/register/?username=%v&password=%v", username, password), "", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)
}

// 重复注册
func TestRegisterAgain(t *testing.T) {
	username := "testAgain"
	password := "123456"

	// 检查是否被注册
	var count int64
	err := repository.DB.Table("user").Where("user_name = ?", username).Count(&count).Error
	assert.Equal(t, err, nil)

	// 没注册则进行注册
	if count == 0 {
		err := repository.DB.Table("user").Create(&repository.User{UserName: username, Password: password}).Error
		assert.Equal(t, err, nil)
	}

	// 注册
	response, err := http.Post(fmt.Sprintf("http://127.0.0.1:8080/douyin/user/register/?username=%v&password=%v", username, password), "", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"].(string), "该用户名已被注册")
}

// 测试已注册用户登录
func TestLogin(t *testing.T) {
	username := "testLogin"
	password := "123456"

	// 检查是否被注册
	var count int64
	err := repository.DB.Table("user").Where("user_name = ?", username).Count(&count).Error
	assert.Equal(t, err, nil)

	// 没注册则进行注册
	if count == 0 {
		err := repository.DB.Table("user").Create(&repository.User{UserName: username, Password: password}).Error
		assert.Equal(t, err, nil)
	}

	// 登录
	response, err := http.Post(fmt.Sprintf("http://127.0.0.1:8080/douyin/user/login/?username=%v&password=%v", username, password), "", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)

}

// 测试密码错误登录
func TestLoginWithWrongPassword(t *testing.T) {
	username := "testLogin"
	password := "123456"

	// 检查是否被注册
	var count int64
	err := repository.DB.Table("user").Where("user_name = ?", username).Count(&count).Error
	assert.Equal(t, err, nil)

	// 没注册则进行注册
	if count == 0 {
		err := repository.DB.Table("user").Create(&repository.User{UserName: username, Password: password}).Error
		assert.Equal(t, err, nil)
	} else {
		// 获得password
		err := repository.DB.Table("user").Where("user_name = ?", username).Update("password", password).Error
		assert.Equal(t, err, nil)
	}

	// 登录
	// 拼接时密码+123
	response, err := http.Post(fmt.Sprintf("http://127.0.0.1:8080/douyin/user/login/?username=%v&password=%v123", username, password), "", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"].(string), "用户密码错误")

}

// 测试未注册用户登录
func TestLoginWithouRegister(t *testing.T) {
	username := "testLogin"
	password := "123456"

	// 检查是否被注册
	var count int64
	err := repository.DB.Table("user").Where("user_name = ?", username).Count(&count).Error
	assert.Equal(t, err, nil)

	// 注册过则删除
	if count == 1 {
		err := repository.DB.Table("user").Where("user_name = ?", username).Delete(&repository.User{UserName: username}).Error
		assert.Equal(t, err, nil)
	}

	// 登录
	response, err := http.Post(fmt.Sprintf("http://127.0.0.1:8080/douyin/user/login/?username=%v&password=%v", username, password), "", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil)

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"].(string), "用户名不存在,请注册")

}

//

// 测试关注操作

// 测试关注列表

// 测试粉丝列表

// --------------------benchmark----------------------------------------

func BenchmarkFeed(b *testing.B) {
	b.SetParallelism(10)

	req := fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v", time.Now().UnixMilli())

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			http.Get(req)
		}
	})
}

// func BenchmarkRegister(b *testing.B) {

// }
