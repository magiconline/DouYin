package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting server")
	r := gin.Default()
	r.GET("/", nil)

	r.Run()
}
