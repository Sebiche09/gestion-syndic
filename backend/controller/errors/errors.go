package errors

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error, message string, statusCode int) {
	if err != nil {
		wrappedErr := fmt.Errorf("%s: %w", message, err)
		fmt.Printf("Error: %v\n", wrappedErr)
		c.JSON(statusCode, gin.H{"error": wrappedErr.Error()})
	} else {
		c.JSON(statusCode, gin.H{"error": message})
	}
}
