package routes

import (
	"github.com/Sebiche09/gestion-syndic/controller"
	"github.com/gin-gonic/gin"
)

func OccupantRoute(router *gin.Engine) {
	router.GET("/", controller.GetOccupant)
	router.POST("/", controller.CreateOccupant)
	router.DELETE("/:id", controller.DeleteOccupant)
	router.PUT("/:id", controller.UpdateOccupant)
}
