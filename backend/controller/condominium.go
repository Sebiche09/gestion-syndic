package controller

import (
	"net/http"

	"github.com/Sebiche09/gestion-syndic/config"
	"github.com/Sebiche09/gestion-syndic/models"
	"github.com/gin-gonic/gin"
)

func CreateCondominium(c *gin.Context) {
	// Parse the request body into the Condominium struct
	var condominium models.Condominium
	if err := c.ShouldBindJSON(&condominium); err != nil {
		handleError(c, err, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Create the condominium in the database
	if err := config.DB.Create(&condominium).Error; err != nil {
		handleError(c, err, "Error creating condominium", http.StatusInternalServerError)
		return
	}

	// Return the created condominium
	c.JSON(http.StatusCreated, condominium)
}
