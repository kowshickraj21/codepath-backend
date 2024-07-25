package main

import (
	"main/models"
	"strconv"

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

	router.GET("/user/:username", func(ctx *gin.Context) {
		username := ctx.Param("username");
		user,err := controllers.GetUser(db,username);

		if err != nil {
			ctx.JSON(500,err);
		}

		ctx.JSON(200,user)
	})

	router.GET("/problem/:problemId", func(ctx *gin.Context) {
		idstr := ctx.Param("problemId");

		id,err := strconv.Atoi(idstr)

		if err != nil {
			ctx.JSON(500, gin.H{"error": err})
		}

		controllers.ViewProblem(db,id)
	})

	router.GET("/problems", func(ctx *gin.Context) {
		
		problems,err := controllers.FetchProblems(db);
		if err != nil {
			ctx.JSON(500,err);
		}
		ctx.JSON(200,problems);
	})

	router.Run(":3050")
}