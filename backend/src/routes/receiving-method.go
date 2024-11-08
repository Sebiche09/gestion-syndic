package routes

import (
	"github.com/Sebiche09/gestion-syndic/src/controller"
	"github.com/gin-gonic/gin"
)

func ReceivingMethodRoute(router *gin.Engine) {
	router.GET("/getdocumentreceivingmethods", controller.GetDocumentReceivingMethods)
	router.GET("/getreminderreceivingmethods", controller.GetReminderReceivingMethods)
}
