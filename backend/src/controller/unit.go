package controller

import (
	"net/http"

	"github.com/Sebiche09/gestion-syndic/src/config"
	"github.com/Sebiche09/gestion-syndic/src/models"

	"github.com/gin-gonic/gin"
)

// GetUnits retourne la liste des unités avec leurs informations associées
// @Summary Liste toutes les unités
// @Description Récupère la liste des unités avec leurs références cadastrales, étages, descriptions, quotas, et les détails des condominiums, adresses, et types d'unités associés
// @Tags Units
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "Liste des unités"
// @Failure 500 {object} map[string]string "Erreur serveur"
// @Router /units [get]
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
