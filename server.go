package main

import (
	"DouYin/controller"
	"DouYin/logger"
	"DouYin/repository"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
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
	}

	r := gin.Default()

	// 托管静态资源
	r.Static("/static", "./static")

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
	//r.GET("/douyin/favorite/action/", func(ctx *gin.Context) {
	//	body := controller.Favorite(ctx)
	//	ctx.JSON(200, body)
	//})
	logger.Logger.Println("启动服务器")
	r.Run()
}
