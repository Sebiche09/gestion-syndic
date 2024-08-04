package routes

import (
	"github.com/Sebiche09/gestion-syndic/controller"
	"github.com/gin-gonic/gin"
)

func InvoiceRoute(router *gin.Engine) {
	router.POST("/upload", controller.UploadHandler)
}
