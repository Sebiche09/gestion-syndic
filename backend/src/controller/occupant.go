package controller

import (
	"github.com/Sebiche09/gestion-syndic/src/config"
	"github.com/Sebiche09/gestion-syndic/src/models"
	"github.com/gin-gonic/gin"
)

func GetOccupant(c *gin.Context) {
	occupants := []models.Occupant{}
	config.DB.Find(&occupants)
	c.JSON(200, &occupants)
}
func CreateOccupant(c *gin.Context) {
	var occupant models.Occupant
	c.BindJSON(&occupant)
	config.DB.Create(&occupant)
	c.JSON(200, &occupant)
}
func DeleteOccupant(c *gin.Context) {
	var occupant models.Occupant
	config.DB.Where("id = ?", c.Param("id")).Delete(&occupant)
	c.JSON(200, &occupant)
}

func UpdateOccupant(c *gin.Context) {
	var occupant models.Occupant
	config.DB.Where("id = ?", c.Param("id")).First(&occupant)
	c.BindJSON(&occupant)
	config.DB.Save(&occupant)
	c.JSON(200, &occupant)
}
