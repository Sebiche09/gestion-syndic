package controller

import (
	"github.com/Sebiche09/gestion-syndic/src/config"
	"github.com/Sebiche09/gestion-syndic/src/models"
	"github.com/gin-gonic/gin"
)

// @Summary Récupère la liste des occupants
// @Description Retourne tous les occupants enregistrés
// @Tags Occupants
// @Produce json
// @Success 200 {array} models.Occupant
// @Router /api/occupants [get]
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
