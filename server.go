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
	} else {
		fmt.Println("数据库连接成功")
	}

	// test
	// videoTable := repository.VideoTable{2, 0, "test", "test", 0, 0, 0}
	// err = repository.InsertVideoTable(&videoTable)
	// if err != nil {
	// 	logger.Logger.Println("数据插入失败,", err)
	// } else {
	// 	logger.Logger.Println("数据插入成功")
	// }

	r := gin.Default()

	// r.GET("/douyin/user", func(ctx *gin.Context) {
	// 	ctx.JSON(json.SUCCESS, gin.H{
	// 		"Code": json.SUCCESS,
	// 		"Msg":  json.GetMsg(json.SUCCESS),
	// 		"Data": "返回数据成功",
	// 	})
	// })
	// r.GET("/douyin/user/login",func(ctx *gin.Context) {
	// 	body  := controller.
	// } )
	// r.GET("/douyin/user/register",)
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
