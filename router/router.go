package router

import (
	"github.com/gin-gonic/gin"
	"pandaBook/controller"
)

func Init() *gin.Engine {
	router := gin.Default()

	router.GET("/d4j/list", func(c *gin.Context) {
		new(controller.D4jController).List(c)
	})

	router.GET("/it/list", func(c *gin.Context) {
		new(controller.ItPandaController).List()
	})

	return router
}
