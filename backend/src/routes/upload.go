package routes

import (
	"github.com/Sebiche09/gestion-syndic/src/controller/upload"
	"github.com/gin-gonic/gin"
)

func InvoiceRoute(router *gin.Engine) {
	router.POST("/upload", upload.UploadHandler)
}
