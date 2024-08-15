package main

import (
	"database/sql"
	"fmt"
	"main/auth"
	"main/controllers"
	"main/initializers"
	"main/models"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{}
var db = &sql.DB{}

func init(){
	initializers.LoadEnv();
	db = initializers.ConnectDB();

	googleOauthConfig = &oauth2.Config{
		ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL: os.Getenv("REDIRECT_URL"),
		Scopes: []string{"profile","email"},
		Endpoint: google.Endpoint,
	}
}


func main() {

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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

	router.POST("/submit/:problemId",func(ctx *gin.Context) {
		id := ctx.Param("problemId");

		var code models.Code
		ctx.ShouldBindJSON(&code)
		fmt.Println(code.Language)
		token,err:= controllers.CreateReq(db,code,id);
		if err != nil {
			ctx.JSON(500,err);
		}

		var res *models.Judge0Response
		time.Sleep(1 * time.Second)
		res,err = controllers.GetReq(token);
		if err != nil {
			ctx.JSON(500,err);
		}
		ctx.JSON(200,res);
	})

	router.GET("/code/:problemId/:language",func(ctx *gin.Context) {
		id := ctx.Param("problemId");
		language := ctx.Param("language");

		fileurl := "problems/$1/boilerplate.$2.txt"
		fileurl = strings.Replace(fileurl,"$1",id,1)
		fileurl = strings.Replace(fileurl,"$2",language,1)
		
		file,err := os.ReadFile(fileurl)
		if err != nil {
			ctx.JSON(500,err);
		}
		boilerplate := string(file)
		fmt.Println(boilerplate)
		ctx.JSON(200,boilerplate);
	})

	router.GET("/auth/google/callback", func(ctx *gin.Context) {
		code := ctx.Query("code")

		user,err := auth.GetUserInfo(code)
		if err != nil {
			ctx.JSON(500,err)
		}

		signedToken,err := auth.SignJWT(user)
		if err != nil {
			ctx.JSON(500,err)
		}

		ctx.JSON(200,signedToken)
		 
	})

	router.Run(":3050")
}