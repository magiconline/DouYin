package main

import (
	"DouYin/controller"
	"DouYin/repository"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	err := repository.Init()
	if err != nil {
		fmt.Println("数据库连接错误:", err)
		os.Exit(1)
	}

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

	r.Run()
}
