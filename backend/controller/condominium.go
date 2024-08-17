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

func GetCivilities(c *gin.Context) {
	var civilities []models.Civility
	if err := config.DB.Find(&civilities).Error; err != nil {
		handleError(c, err, "Error fetching civilities", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, civilities)
}

func GetReceivingMethods(c *gin.Context) {
	var receivingMethods []models.ReceivingMethod
	if err := config.DB.Find(&receivingMethods).Error; err != nil {
		handleError(c, err, "Error fetching receiving methods", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, receivingMethods)
}
