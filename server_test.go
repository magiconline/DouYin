package main

import (
	"DouYin/repository"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var r *gin.Engine

func TestMain(m *testing.M) {
	// 初始化数据库连接
	err := repository.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 初始化服务器
	gin.SetMode(gin.ReleaseMode)
	file, _ := os.Open("log")
	gin.DefaultWriter = file
	r = setupRouter()

	code := m.Run()

	file.Close()
	os.Exit(code)
}

func TestFeed(t *testing.T) {
	// 初始化请求
	timeStamp := time.Now().UnixMilli()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/douyin/feed/?latest_time=%v", timeStamp), nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func BenchmarkFeed(b *testing.B) {
	// 初始化请求
	timeStamp := time.Now().UnixMilli()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/douyin/feed/?latest_time=%v", timeStamp), nil)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.ServeHTTP(w, req)
		}

	})
}
