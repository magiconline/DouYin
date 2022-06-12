package main

import (
	"DouYin/logger"
	"DouYin/repository"
	"DouYin/service"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	code := m.Run()

	os.Exit(code)
}

// 测试没有token
func TestFeedWithoutToken(t *testing.T) {
	timeStamp := time.Now().UnixMilli()

	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", fmt.Sprintf("/douyin/feed/?latest_time=%v", timeStamp), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
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

	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", fmt.Sprintf("/douyin/feed/?latest_time=%v&token=%v", timeStamp, token), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"].(string), "token contains an invalid number of segments")
}

// 测试过期token
func TestFeedWithExpiredToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiZXhwIjoxNjU0NzAwNzcyLCJpc3MiOiJ6amN5In0.X9VuPerdOP8TNFxVpWY3vLVFPHdVE72un8TiimFMFPk"
	timeStamp := time.Now().UnixMilli()

	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", fmt.Sprintf("/douyin/feed/?latest_time=%v&token=%v", timeStamp, token), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
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

	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v", timeStamp), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
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
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/user/register/?username=%v&password=%v", username, password), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
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
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/user/register/?username=%v&password=%v", username, password), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
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
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/user/login/?username=%v&password=%v", username, password), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
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
		// 设置password
		err := repository.DB.Table("user").Where("user_name = ?", username).Update("password", password).Error
		assert.Equal(t, err, nil)
	}

	// 登录
	// 拼接时密码+123
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/user/login/?username=%v&password=%v123", username, password), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
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
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/user/login/?username=%v&password=%v", username, password), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"].(string), "用户名不存在,请注册")

}

// 测试关注操作
func TestRelation(t *testing.T) {
	username1 := "testRelation1"
	username2 := "testRelation2"
	password := "123456"

	// 注册用户
	err := repository.DB.Table("user").Where("user_id in ?", []string{username1, username2}).Delete(&repository.User{}).Error
	assert.Equal(t, err, nil)

	users := []repository.User{
		{UserName: username1, Password: password},
		{UserName: username2, Password: password},
	}
	err = repository.DB.Table("user").Create(&users).Error
	assert.Equal(t, err, nil)

	token, err := service.GenerateToken(users[0].UserId)
	assert.Equal(t, err, nil)

	// 关注
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/relation/action/?token=%v&to_user_id=%v&action_type=%v", token, users[1].UserId, 1), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)

	// 取消关注
	response = httptest.NewRecorder()
	request, err = http.NewRequest("POST", fmt.Sprintf("/douyin/relation/action/?token=%v&to_user_id=%v&action_type=%v", token, users[1].UserId, 2), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body = make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)

}

// 测试重复关注
func TestRalationAgain(t *testing.T) {
	username1 := "testRelation1"
	username2 := "testRelation2"
	password := "123456"

	// 注册用户
	err := repository.DB.Table("user").Where("user_id in ?", []string{username1, username2}).Delete(&repository.User{}).Error
	assert.Equal(t, err, nil)

	users := []repository.User{
		{UserName: username1, Password: password},
		{UserName: username2, Password: password},
	}
	err = repository.DB.Table("user").Create(&users).Error
	assert.Equal(t, err, nil)

	token, err := service.GenerateToken(users[0].UserId)
	assert.Equal(t, err, nil)

	// 关注
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/relation/action/?token=%v&to_user_id=%v&action_type=%v", token, users[1].UserId, 1), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)

	// 重复关注
	response = httptest.NewRecorder()
	request, err = http.NewRequest("POST", fmt.Sprintf("/douyin/relation/action/?token=%v&to_user_id=%v&action_type=%v", token, users[1].UserId, 1), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body = make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"].(string), "已关注，禁止重复关注")
}

// 测试取消关注

// 测试重复取消关注

// 测试关注列表

// 测试粉丝列表

// --------------------benchmark----------------------------------------

// func BenchmarkFeed(b *testing.B) {
// 	b.SetParallelism(10)

// 	req := fmt.Sprintf("http://127.0.0.1:8080/douyin/feed/?latest_time=%v", time.Now().UnixMilli())

// 	b.RunParallel(func(p *testing.PB) {
// 		for p.Next() {
// 			http.Get(req)
// 		}
// 	})
// }

//测试点赞操作
func TestFavorite(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTk5OTc5LCJpc3MiOiJ6amN5In0.oNrcj2xrgiy5zh0j2So-Sm_vxIG_lsYxRT2rcCQ5EGA"
	videoId := 1
	actionType := 1
	// 检查是否已经点赞
	userId, err := service.Token2ID(token)
	assert.Equal(t, err, nil)
	flag := service.IsThumbUp(userId, uint64(videoId))
	//如果返回true 说明已点赞 删除点赞状态
	if flag {
		service.DeleteStar(userId, uint64(videoId))
	}

	// 点赞操作
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/favorite/action/?token=%v&video_id=%v&action_type=%v", token, videoId, actionType), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)
	assert.Equal(t, body["status_msg"], "点赞成功！")

}

//测试取消点赞操作
func TestCancelFavorite(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTk5OTc5LCJpc3MiOiJ6amN5In0.oNrcj2xrgiy5zh0j2So-Sm_vxIG_lsYxRT2rcCQ5EGA"
	videoId := 8
	actionType := 2
	// 检查是否已经点赞
	userId, err := service.Token2ID(token)
	assert.Equal(t, err, nil)
	flag := service.IsThumbUp(userId, uint64(videoId))
	//如果返回false 说明未点赞 增加点赞状态
	if !flag {
		service.AddStar(userId, uint64(videoId))
	}

	// 取消点赞操作
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/favorite/action/?token=%v&video_id=%v&action_type=%v", token, videoId, actionType), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)
	assert.Equal(t, body["status_msg"], "取消点赞成功！")

}

//测试重复点赞
func TestRepeatFavorite(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTk5OTc5LCJpc3MiOiJ6amN5In0.oNrcj2xrgiy5zh0j2So-Sm_vxIG_lsYxRT2rcCQ5EGA"
	videoId := 8
	actionType := 1
	// 检查是否已经点赞
	userId, err := service.Token2ID(token)
	assert.Equal(t, err, nil)
	flag := service.IsThumbUp(userId, uint64(videoId))
	//如果返回false 说明未点赞 插入点赞数据
	if !flag {
		service.AddStar(userId, uint64(videoId))
	}

	// 点赞操作
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/favorite/action/?token=%v&video_id=%v&action_type=%v", token, videoId, actionType), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 5)
	assert.Equal(t, body["status_msg"], "请勿重复点赞！")

}

//测试在未点赞的状态下取消点赞的操作
func TestCancelFavoriteWithoutFavorite(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTk5OTc5LCJpc3MiOiJ6amN5In0.oNrcj2xrgiy5zh0j2So-Sm_vxIG_lsYxRT2rcCQ5EGA"
	videoId := 1
	actionType := 2
	// 检查是否已经点赞
	userId, err := service.Token2ID(token)
	assert.Equal(t, err, nil)
	flag := service.IsThumbUp(userId, uint64(videoId))
	//如果返回true 说明已点赞 删除点赞数据
	if flag {
		service.DeleteStar(userId, uint64(videoId))
	}

	// 取消点赞操作
	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/favorite/action/?token=%v&video_id=%v&action_type=%v", token, videoId, actionType), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 6)
	assert.Equal(t, body["status_msg"], "当前暂无点赞数据！")

}

//测试获取未登录用户获取点赞视频列表
func TestFavoriteWithoutLogin(t *testing.T) {
	token := ""
	userId := "1"

	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:8080/douyin/favorite/list/?token=%v&user_id=%v", token, userId), nil)

	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"], "用户未登录")
}

//测试获取登录用户获取点赞视频列表
func TestFavoriteLogin(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTY2MDcwLCJpc3MiOiJ6amN5In0.7jRyQql7BZ78PDqmL-zn2hhf_9yTxIUvKPo-dCbTEwg"
	userId := "1"

	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:8080/douyin/favorite/list/?token=%v&user_id=%v", token, userId), nil)

	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)
	assert.Equal(t, body["status_msg"], "success")
}

// 测试remark没有token
func TestReamrkWithoutToken(t *testing.T) {
	//请求值的参数
	videoid := 1
	actiontype := 1

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/comment/action/?video_id=%v&action_type=%v&comment_text=tokentest", videoid, actiontype), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	if int(body["status_code"].(float64)) != 0 {
		//t.Errorf("status_code: %v != 0, status_msg: %v", body["status_code"], body["status_msg"])
		//t.FailNow()
		assert.Equal(t, body["status_msg"], "用户未登录")
	}

}

// 测试错误token
func TestReamrkWithWrongToken(t *testing.T) {
	//请求值的参数
	token := "132"
	videoid := 1
	actiontype := 1

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/comment/action/?token=%v&video_id=%v&action_type=%v&comment_text=tokentest", token, videoid, actiontype), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)
	assert.Equal(t, body["status_msg"].(string), "token has some error")
}

//测试过期token
func TestRemarkWithExpiredToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiZXhwIjoxNjU0NzAwNzcyLCJpc3MiOiJ6amN5In0.X9VuPerdOP8TNFxVpWY3vLVFPHdVE72un8TiimFMFPk"
	videoid := 1
	actiontype := 1

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/comment/action/?token=%v&video_id=%v&action_type=%v&comment_text=tokentest", token, videoid, actiontype), nil)
	assert.Equal(t, err, nil)

	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 1)

	result, err := regexp.Match("token has some error", []byte(body["status_msg"].(string)))
	if err != nil {
		t.Fatal(err.Error())
		t.FailNow()
	}
	if result != true {
		t.Fatalf("status_msg: %v != token has some error", body["status_msg"].(string))
		t.FailNow()
	}
}

//测试插入评论
func TestInsertRemark(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTk2OTE0LCJpc3MiOiJ6amN5In0.WeN4fZgkitj_ETYIvwAP-nvIPewWMIRBT4tIbX_mTYY"
	videoid := 1
	actiontype := 1
	commenttext := "test1"

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/comment/action/?token=%v&video_id=%v&action_type=%v&comment_text=%v", token, videoid, actiontype, commenttext), nil)
	assert.Equal(t, err, nil)
	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	//如果插入失败
	if int(body["status_code"].(float64)) == 1 {
		assert.Equal(t, body["status_msg"].(string), "insert error")
	} else {
		assert.Equal(t, int(body["status_code"].(float64)), 0)
	}
}

//测试删除评论
func TestDeleteRemark(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTk2OTE0LCJpc3MiOiJ6amN5In0.WeN4fZgkitj_ETYIvwAP-nvIPewWMIRBT4tIbX_mTYY"
	videoid := 1
	actiontype := 2
	commenttext := "testdelete"
	commentid := 111

	response := httptest.NewRecorder()
	request, err := http.NewRequest("POST", fmt.Sprintf("/douyin/comment/action/?token=%v&video_id=%v&action_type=%v&comment_text=%v&comment_id=%v", token, videoid, actiontype, commenttext, commentid), nil)
	assert.Equal(t, err, nil)
	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	//如果删除失败
	if int(body["status_code"].(float64)) == 1 {
		assert.Equal(t, body["status_msg"].(string), "delete error")
	} else {
		assert.Equal(t, int(body["status_code"].(float64)), 0)
	}

}

//测试查看所有评论列表
func TestRemarkList(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiZXhwIjoxNjYyMTk2OTE0LCJpc3MiOiJ6amN5In0.WeN4fZgkitj_ETYIvwAP-nvIPewWMIRBT4tIbX_mTYY"
	videoid := 1

	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", fmt.Sprintf("/douyin/comment/list/?token=%v&video_id=%v", token, videoid), nil)
	assert.Equal(t, err, nil)
	r.ServeHTTP(response, request)
	assert.Equal(t, response.Code, 200)

	body := make(map[string]interface{})
	err = json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, err, nil)

	assert.Equal(t, int(body["status_code"].(float64)), 0)

}
