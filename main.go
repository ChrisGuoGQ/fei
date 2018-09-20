package main

import (
	"gintest/controllers"
	"gintest/models"
	"gintest/utils"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
func main() {
	models.Bootstrap()
	utils.InitConfig()
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/groups", controllers.IndexGroup)
	r.POST("/groups", controllers.CreateGroup)
	r.PUT("/groups/:id", controllers.UpdateGroup)
	r.DELETE("/groups/:id", controllers.DestroyGroup)
	r.GET("/cameras", controllers.IndexCamera)
	r.POST("/cameras", controllers.CreateCamera)
	r.PUT("/cameras/:id", controllers.UpdateCamera)
	r.DELETE("/cameras/:id", controllers.DestroyCamera)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run() // listen and serve on 0.0.0.0:8080
}
