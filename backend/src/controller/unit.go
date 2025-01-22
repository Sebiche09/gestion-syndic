package controller

import (
	"net/http"

	"github.com/Sebiche09/gestion-syndic/src/config"
	"github.com/Sebiche09/gestion-syndic/src/models"

	"github.com/gin-gonic/gin"
)

func GetUnits(c *gin.Context) {
	db := config.DB

	var units []models.Unit

	if err := db.Preload("Condominium").
		Preload("Address").
		Preload("UnitType").
		Find(&units).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "test"})
		return
	}

	response := []gin.H{}
	for _, unit := range units {
		response = append(response, gin.H{
			"id":                  unit.ID,
			"cadastral_reference": unit.CadastralReference,
			"floor":               unit.Floor,
			"description":         unit.Description,
			"quota":               unit.Quota,
			"condominium": gin.H{
				"id":   unit.Condominium.ID,
				"name": unit.Condominium.Name,
			},
			"address": gin.H{
				"city": unit.Address.City,
			},
			"unit_type": gin.H{
				"id":    unit.UnitType.ID,
				"label": unit.UnitType.Label,
			},
		})
	}
	c.JSON(http.StatusOK, response)
}
