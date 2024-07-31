package main

import (
	"fmt"
	"main/controllers"
	"main/initializers"
	"main/models"
	"strconv"
	"time"

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

		problem,err := controllers.ViewProblem(db,id);
		if err != nil {
			ctx.JSON(500,err);
		}
		ctx.JSON(200,problem);
	})

	router.GET("/problems", func(ctx *gin.Context) {
		
		problems,err := controllers.FetchProblems(db);
		if err != nil {
			ctx.JSON(500,err);
		}
		ctx.JSON(200,problems);
	})

	router.POST("/submit",func(ctx *gin.Context) {
		var code models.Code
		ctx.ShouldBindJSON(&code)
		token,err:= controllers.CreateReq(code.Code);
		if err != nil {
			ctx.JSON(500,err);
		}
		fmt.Println(token)
		var res *models.Judge0Response
		time.Sleep(1 * time.Second)
		res,err = controllers.GetReq(token);
		if err != nil {
			ctx.JSON(500,err);
		}
		ctx.JSON(200,res);
	})
	router.Run(":3050")
}