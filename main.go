package main

import (
	"kukuh/go-gin-library-project/app/controller"
	"kukuh/go-gin-library-project/app/repository"
	"kukuh/go-gin-library-project/app/service"
	"kukuh/go-gin-library-project/database"
	"kukuh/go-gin-library-project/helper/token"
	"kukuh/go-gin-library-project/response"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {

	db, err := database.NewMysqlClient()
	if err != nil {
		panic(err)
	}

	// Initialize validator
	validate := validator.New()

	// Initialize repositories
	userRepository := repository.NewUserRepository()
	bookRepository := repository.NewBookRepository()
	borrowingRepository := repository.NewBorrowingRepository()

	// Initialize services
	userService := service.NewUserService(userRepository, db, validate)
	bookService := service.NewBookService(bookRepository, db, validate)
	borrowingService := service.NewBorrowingService(borrowingRepository, bookRepository, db, validate)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	bookController := controller.NewBookController(bookService)
	borrowingController := controller.NewBorrowingController(borrowingService)

	router := gin.Default()

	// API Grouping
	api := router.Group("/api")
	{

		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)

		api.POST("/book", bookController.Create)
		api.PUT("/book/:id", bookController.Update)
		api.GET("/book/:id", bookController.Find)
		api.DELETE("/book/:id", bookController.Delete)
		api.GET("/book", bookController.FindAll)

		auth := api.Group("/auth")
		auth.Use(CheckAuth())
		{
			auth.PUT("/users", userController.UpdateUserOwn)
			auth.DELETE("/users", userController.DeleteUser)
			auth.POST("/borrowing", borrowingController.Create)
			auth.POST("/borrowing/return/:id", borrowingController.Return)
			auth.GET("/borrowing/:id", borrowingController.Find)
			auth.GET("/borrowing", borrowingController.FindAll)
		}
	}

	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.UnauthorizedError("len token must be 2")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		payload, err := token.ValidateJwtToken(bearerToken[1])
		if err != nil {
			resp := response.UnauthorizedError(err.Error())
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}
		ctx.Set("authId", payload.AuthId)
		ctx.Next()
	}
}
