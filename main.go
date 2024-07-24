package main

import (
	"fmt"
	"log"
	"main/controllers"
	"main/initializers"

	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnv();
	db := initializers.ConnectDB();
	userController := controllers.NewUserController(db);
	newUser, err := userController.CreateUser("johndoe", "John Doe", "john.doe@example.com", "http://example.com/picture.jpg", []int{1, 2, 3})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User created with username: %s\n", newUser.Username)
}

func main() {
	router := gin.Default()
	router.GET("/", func(context *gin.Context) {
		context.String(200, "hello")
	})
	router.Run(":3050")
}