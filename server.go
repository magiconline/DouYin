package main

import (
	"DouYin/controller"
	"DouYin/logger"
	"DouYin/repository"
	"DouYin/service"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetOutBoundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "unknown"
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(localAddr.String(), ":")[0]
	return ip
}

func main() {
	// 获得本机ip+port
	ip := GetOutBoundIP()
	port := ":8080"
	addr := "http://" + ip + port
	fmt.Println("start server at ", addr)

	// 设置静态资源ip
	service.SetServerIP(addr)

	// 初始化日志
	err := logger.Init("./log")
	if err != nil {
		fmt.Println("日志初始化失败：", err)
		os.Exit(1)
	}
	logger.Logger.Println("日志初始化成功")

	// 初始化数据库连接
	err = repository.Init()
	if err != nil {
		fmt.Println("数据库连接错误:", err)
		os.Exit(1)
	} else {
		fmt.Println("数据库连接成功")
	}

	r := gin.Default()

	// 托管静态资源
	r.Static("/static", "./static")

	// 路由
	r.POST("/douyin/user/login/", func(ctx *gin.Context) {
		body := controller.UserLogin(ctx)
		ctx.JSON(200, body)
	})
	r.POST("/douyin/user/register/", func(ctx *gin.Context) {
		body := controller.UserRegister(ctx)
		ctx.JSON(200, body)
	})
	r.GET("/douyin/user/", func(ctx *gin.Context) {
		body := controller.UserInfo(ctx)
		ctx.JSON(200, body)
	})
	r.GET("/douyin/feed/", func(ctx *gin.Context) {
		body := controller.Feed(ctx)
		ctx.JSON(200, body)
	})

	r.POST("/douyin/publish/action/", func(ctx *gin.Context) {
		body := controller.PublishAction(ctx)
		ctx.JSON(200, body)
	})

	r.GET("/douyin/publish/list/", func(ctx *gin.Context) {
		body := controller.PublishList(ctx)
		ctx.JSON(200, body)
	})

	r.POST("/douyin/favorite/action/", func(ctx *gin.Context) {
		body := controller.Favorite(ctx)
		ctx.JSON(200, body)
	})
	r.GET("/douyin/favorite/list/", func(ctx *gin.Context) {
		body := controller.ThumbListVideo(ctx)
		ctx.JSON(200, body)
	})
	logger.Logger.Println("启动服务器")
	r.Run(port)
}
