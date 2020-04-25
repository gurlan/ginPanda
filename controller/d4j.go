package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"pandaBook/model"
)

type D4jController struct {
	collector *colly.Collector
}

func (d *D4jController) List(ctx *gin.Context) {
	new(model.D4jModel).List()
}
