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

// 获得本机IP
func getOutBoundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "", err
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(localAddr.String(), ":")[0]
	return ip, nil
}

// 设置静态资源 URL 前缀
func setIP(ip string, port string) {
	addr := "http://" + ip + port
	service.SetServerIP(addr)
	logger.Logger.Println("静态资源URL前缀: ", addr)
}

// 设置路由
func setupRouter(r *gin.Engine) {
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

	//留评论
	r.POST("/douyin/comment/action/", func(ctx *gin.Context) {
		body := controller.Leave_remark(ctx)
		ctx.JSON(200, body)
	})
	//查看视频的所有评论，按发布时间倒序
	r.GET("/douyin/comment/list/", func(ctx *gin.Context) {
		body := controller.View_video_remark(ctx)
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

	r.POST("/douyin/relation/action/", func(ctx *gin.Context) {
		body := controller.RelationAction(ctx)
		ctx.JSON(200, body)
	})
	r.GET("/douyin/relation/follow/list/", func(ctx *gin.Context) {
		body := controller.FollowList(ctx)
		ctx.JSON(200, body)
	})
	r.GET("/douyin/relation/follower/list/", func(ctx *gin.Context) {
		body := controller.FollowerList(ctx)
		ctx.JSON(200, body)
	})
}

func main() {
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
		os.Exit(1)
	}

	port := ":8080"

	setIP(ip, port)

	// 初始化数据库连接
	err = repository.Init()
	if err != nil {
		logger.Fatalln("数据库初始化失败,", err.Error())
		os.Exit(1)
	}

	// 设置release模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化gin日志
	file, err := os.OpenFile("./log/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatalln("gin日志初始化失败,", err.Error())
		os.Exit(1)
	}
	gin.DefaultWriter = file

	// 创建服务器
	r := gin.Default()

	// 设置路由
	setupRouter(r)

	// 托管静态资源
	r.Static("/static", "./static")

	// 启动服务器
	logger.Println("启动服务器")
	fmt.Printf("启动服务器: %v%v\n", ip, port)

	err = r.Run(port)
	if err != nil {
		logger.Logger.Fatal("服务器启动失败,", err.Error())
		os.Exit(1)
	}
}
