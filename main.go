package main

import (
	"main/initializers"

	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnv();
	initializers.ConnectDB();
}

func main() {
	router := gin.Default()
	router.GET("/", func(context *gin.Context) {
		context.String(200, "hello")
	})
	router.Run(":3050")
}