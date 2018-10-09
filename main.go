package main

import (
	"fei/controllers"
	"fei/models"
	"fei/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	utils.InitConfig()
	models.Bootstrap()
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/groups", controllers.IndexGroup)
	r.POST("/groups", controllers.CreateGroup)
	r.PUT("/groups/:id", controllers.UpdateGroup)
	r.DELETE("/groups/:id", controllers.DestroyGroup)
	r.GET("/cameras", controllers.IndexCamera)
	r.POST("/cameras", controllers.CreateCamera)
	r.PUT("/cameras/:id", controllers.UpdateCamera)
	r.DELETE("/cameras/:id", controllers.DestroyCamera)
	r.GET("/cameras/:id/start", controllers.StartCamera)
	r.GET("/cameras/:id/stop", controllers.StopCamera)
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/add_bases", controllers.AddBases)
	r.GET("/sync_bases", controllers.SyncBases)
	cfg := utils.GetCfg()
	if ConnectAll := cfg.ConnectAll; ConnectAll {
		controllers.ConnectAll()
	}
	r.Run(":" + cfg.Port) // listen and serve on 0.0.0.0:8080
}
