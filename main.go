package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	fmt.Println("Hello World")
	router := gin.Default()
	router.GET("/request", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": "Hello World",
		})
	})
	router.Run(":8888")
}
