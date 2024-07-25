package main

import (
	"main/models"

	"main/controllers"
	"main/initializers"

	"github.com/gin-gonic/gin"
)


func main() {
	initializers.LoadEnv();
	db := initializers.ConnectDB();

	router := gin.Default()
	router.POST("/user", func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		controllers.CreateUser(db,user)
	})


	router.Run(":3050")
}