package router

import (
	"github.com/gin-gonic/gin"
	"pandaBook/controller"
)

func Init() *gin.Engine {
	router := gin.Default()

	router.GET("/list", func(c *gin.Context) {
		new(controller.D4jController).List(c)
	})

	router.GET("/detail", func(c *gin.Context) {
		new(controller.D4jController).Detail()
	})


	return router
}
