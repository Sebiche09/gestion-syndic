package routes

import (
	"github.com/Sebiche09/gestion-syndic/controller"
	"github.com/gin-gonic/gin"
)

func CivilityRoute(router *gin.Engine) {
	router.GET("/getcivilities", controller.GetCivilities)
}
