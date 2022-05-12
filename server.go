package main

import (
	"DouYin/controller"
	"DouYin/repository"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	err := logger.Init("./log")
	if err != nil {
		fmt.Println("日志初始化错误", err)
		os.Exit(1)
	}
	logger.Logger.Println("日志初始化成功")

	// 初始化数据库连接
	err = repository.Init()
	if err != nil {
		logger.Logger.Println("数据库初始化错误:", err)
		os.Exit(1)
	}
	logger.Logger.Println("数据库初始化成功")

	fmt.Println("Starting server")
	r := gin.Default()

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

	logger.Logger.Println("启动服务器")
	r.Run()
}
