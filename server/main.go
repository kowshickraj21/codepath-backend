package main

import (
	"database/sql"
	"fmt"
	"main/aws"
	"main/controllers"
	"main/initializers"
	"main/models"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{}
var db = &sql.DB{}
var client *s3.Client

func init(){
	initializers.LoadEnv();
	db = initializers.ConnectDB();
	client = initializers.AWSInit();

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
		AllowOrigins:     []string{os.Getenv("FRONTEND_ORIGIN")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization","user"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/user",func(ctx *gin.Context) {
	 jwt := ctx.Query("token")

	 user,err := controllers.GetAuthUser(db,jwt)
	 if(err != nil){
		ctx.JSON(500,err);
		return;
	 }
	 ctx.JSON(200,user);
	 
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

	router.GET("/submission/:problemId", func(ctx *gin.Context) {
		idstr := ctx.Param("problemId");
		id,err := strconv.Atoi(idstr)
		if err != nil {
			ctx.JSON(500,err);
			return
		}
		jwt := ctx.Request.Header.Get("user")

		solutions,err := controllers.HandleSolutions(db,id,jwt)
		ctx.JSON(200,solutions)
	})

	router.GET("/problems", func(ctx *gin.Context) {
		
		problems,err := controllers.FetchProblems(db);
		if err != nil {
			ctx.JSON(500,err);
			return
		}
		ctx.JSON(200,problems);
	})

	router.POST("/submit/:problemId",func(ctx *gin.Context) {
		id := ctx.Param("problemId");
		jwt := ctx.Request.Header.Get("user")

		var code models.Code
		ctx.ShouldBindJSON(&code)
		res,err:= controllers.HandleSubmissions(db,client,code,id,jwt);
		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}
		ctx.JSON(200,res);
	})

	router.POST("/run/:problemId",func(ctx *gin.Context) {
		id := ctx.Param("problemId");
		jwt := ctx.Request.Header.Get("user")

		var code models.Code
		ctx.ShouldBindJSON(&code)
		res,err:= controllers.HandleRun(db,client,code,id,jwt);
		if err != nil {
			ctx.JSON(500,err.Error());
			return
		}
		ctx.JSON(200,res);
	})
	
	router.GET("/code/:problemId/:language",func(ctx *gin.Context) {
		id := ctx.Param("problemId");
		language := ctx.Param("language");

		fileurl := "$1/boilerplate.$2.txt"
		fileurl = strings.Replace(fileurl,"$1",id,1)
		fileurl = strings.Replace(fileurl,"$2",language,1)
		
		boilerplate := aws.ReadFile(context.TODO(),client,os.Getenv("AWS_BUCKET"),fileurl)
		fmt.Println(boilerplate)
		ctx.JSON(200,boilerplate);
	})

	router.GET("/auth/google/callback", func(ctx *gin.Context) {
		code := ctx.Query("code")

		token,err := controllers.HandleGoogleUser(db,code)
		if err != nil {
			ctx.JSON(500,err)
			return
		}

		ctx.JSON(200,token)
		 
	})

	router.GET("/auth/github/callback", func(ctx *gin.Context) {
		code := ctx.Query("code")
		token,err := controllers.HandleGithubUser(db,code) 
		fmt.Println("err: ",err)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200,token)
	})

	router.Run(fmt.Sprintf(":%s",os.Getenv("BACKEND_PORT")))
}