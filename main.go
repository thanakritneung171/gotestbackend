package main

import (
	"gotestbackend/controllers"
	"gotestbackend/database"
	"gotestbackend/middlewares"

	//"gotestbackend/middlewares"

	_ "gotestbackend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	db = database.SetupDB() // Initialize your database connection
)

//	@title			Thanakrit GOlang test Rest API
//	@version		1.0
//	@description	This is a sample server for a Gin REST API. \n Authorize  use Bearer {token from login} \n example  "Bearer eyJhbGciOiJIUzI"

// @host		localhost:8080
// @BasePath	/api

// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @in header
// @name Authorization

func main() {

	// Run migrations
	database.Migrate(db)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Routes
	v := r.Group("/api")
	{
		//CRUD
		v.GET("/userAll", controllers.GetAllUser)
		//5.
		v.POST("/user/register", controllers.Register)
		v.GET("/user/GetUserByID/:id", controllers.GetUserByID)
		v.PUT("/user/UpdateUserByID/:id", controllers.UpdateUserByID)
		v.DELETE("/user/DeleteUserByID/:id", controllers.DeleteUserByID)
		//CRUD
		//6.
		v.POST("/user/login", controllers.Login) //
		//v1.Use(middlewares.AuthMiddleware())
	}
	v1 := r.Group("/api").Use(middlewares.JWTAuthMiddleware())
	{
		//7.
		v1.GET("/user/me", controllers.GetUser) //
		//8.
		v1.PATCH("/user/me", controllers.UpdateUser)
		//9.
		v1.POST("/accounting/transfer", controllers.Transfer)
		//10.
		v1.GET("/accounting/transfer-list", controllers.GetTransferList)
	}

	// Swagger route
	r.StaticFile("/swagger.json", "./docs/swagger.json")

	r.Run(":8080")
}
