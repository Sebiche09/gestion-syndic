package upload

import (
	"log"

	"github.com/gin-gonic/gin"
)

// handleError centralise la gestion des erreurs
func handleError(c *gin.Context, err error, message string, statusCode int) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
	c.JSON(statusCode, gin.H{"error": message})
}
