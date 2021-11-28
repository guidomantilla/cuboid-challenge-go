package router

import (
	"cuboid-challenge/app/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup returns the app router.
func Setup() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	bag := router.Group("/bags")
	{
		bag.GET("", controller.ListBags)
		bag.GET("/:bagID", controller.GetBag)
		bag.POST("", controller.CreateBag)
		bag.DELETE("/:bagID", controller.DeleteBag)
	}

	cuboid := router.Group("/cuboids")
	{
		cuboid.GET("", controller.ListCuboids)
		cuboid.GET("/:cuboidID", controller.GetCuboid)
		cuboid.POST("", controller.CreateCuboid)
		cuboid.PUT("/:cuboidID", controller.UpdateCuboid)
		cuboid.DELETE("/:cuboidID", controller.DeleteCuboid)
	}

	return router
}
