package routes

import (
	"github.com/Sebiche09/gestion-syndic/controller"
	"github.com/gin-gonic/gin"
)

func CondominiumRoute(router *gin.Engine) {
	router.POST("/condominium", controller.CreateCondominium)
}
