package routes

import (
	"github.com/Sebiche09/gestion-syndic/src/controller"
	"github.com/gin-gonic/gin"
)

func UnitRoute(router *gin.Engine) {
	router.GET("/unit", controller.GetUnits)
}
