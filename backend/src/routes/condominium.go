package routes

import (
	"github.com/Sebiche09/gestion-syndic/src/controller"
	"github.com/gin-gonic/gin"
)

func CondominiumRoute(router *gin.Engine) {
	router.POST("/condominium", controller.CreateCondominium)
	router.GET("/check-uniqueness", controller.CheckUniqueness)
	router.GET("/all-condominiums", controller.GetAllCondominiums)
}
