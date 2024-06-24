package main

import (
	"gotestbackend/controllers"
	"gotestbackend/database"
	"gotestbackend/middlewares"

	_ "gotestbackend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	db = database.SetupDB() // Initialize your database connection
)

// @title Thanakrit GOlang test Rest API
// @version 1.0
// @description This is a sample server for a Gin REST API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {

	// Run migrations
	database.Migrate(db)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	v1 := r.Group("/api/v1")
	v1.Use(middlewares.JWTAuthMiddleware())
	{

		//CRUD
		//5.
		v1.POST("/user/register", controllers.Register)
		v1.GET("/userAll", controllers.GetAllUser)
		v1.GET("/user/GetUserByID/:id", controllers.GetUserByID)
		v1.PUT("/user/UpdateUserByID/:id", controllers.UpdateUserByID)
		v1.DELETE("/user/DeleteUserByID/:id", controllers.DeleteUserByID)
		//CRUD

		//7.
		v1.GET("/user/profile/:id", controllers.GetUserProfile)

		// Authenticated routes
		//6.
		v1.POST("/user/login", controllers.Login)
		//v1.Use(middlewares.AuthMiddleware())
		//8.
		v1.PATCH("/user/me", controllers.UpdateUser)
		//9.
		v1.POST("/accounting/transfer", controllers.TransferCredit)
		//10.
		v1.GET("/accounting/transfer-list", controllers.GetTransferList)
	}

	// Swagger route
	r.StaticFile("/swagger.json", "./docs/swagger.json")

	r.Run(":8080")
}
