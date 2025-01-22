package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sebiche09/gestion-syndic/src/config"
	"github.com/Sebiche09/gestion-syndic/src/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CheckIfExists vérifie si une valeur existe déjà dans une table donnée.
// tableName : Nom de la table
// columnName : Nom de la colonne
// value : Valeur à vérifier
func CheckIfExists(db *gorm.DB, tableName string, conditions map[string]interface{}) (bool, error) {
	var count int64
	query := db.Table(tableName).Where(conditions).Count(&count)
	if query.Error != nil {
		return false, query.Error
	}
	return count > 0, nil
}

func getOccupantTypeLabel(title string) (string, error) {
	// Mapping des abréviations aux labels
	abbreviationToLabel := map[string]string{
		"PP":      "pleine propriete",
		"NP":      "nue propriete",
		"US":      "usufruit",
		"SUPERF":  "superficiaire",
		"USA/HAB": "usage/habitation",
		"emph":    "emphyteote",
	}

	// Vérification si le titre contient une des abréviations
	for abbreviation, label := range abbreviationToLabel {
		if strings.Contains(strings.ToUpper(title), abbreviation) {
			return label, nil
		}
	}

	// Retourner une erreur si aucune correspondance n'est trouvée
	return "", fmt.Errorf("no matching occupant type for title: %s", title)
}

func GetAllCondominiums(c *gin.Context) {
	db := config.DB

	var condominiums []models.Condominium

	// Chargement de tous les condominiums avec leurs adresses associées
	if err := db.Preload("Address").Find(&condominiums).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch condominiums"})
		return
	}

	// Construire la réponse avec les champs nécessaires
	response := []gin.H{}
	for _, condo := range condominiums {
		response = append(response, gin.H{
			"name":   condo.Name,
			"prefix": condo.Prefix,
			"city":   condo.Address.City, // Association Address chargée avec Preload
		})
	}
	if len(condominiums) == 0 {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	c.JSON(http.StatusOK, response)
}

func CreateCondominium(c *gin.Context) {
	var requestData struct {
		Informations struct {
			Name        string `json:"name"`
			Prefix      string `json:"prefix"`
			Description string `json:"description"`
		} `json:"informations"`
		Address struct {
			Street            string `json:"street"`
			AddressComplement string `json:"address_complement"`
			City              string `json:"city"`
			PostalCode        string `json:"postal_code"`
			Country           string `json:"country"`
		} `json:"address"`
		Occupants []struct {
			Name                    string `json:"name"`
			Surname                 string `json:"surname"`
			Email                   string `json:"email"`
			Corporation             bool   `json:"corporation"`
			IsConcierge             bool   `json:"isConcierge"`
			Phone                   string `json:"phone"`
			Iban                    string `json:"iban"`
			BirthDate               string `json:"birthdate"`
			Civility                string `json:"civility"`
			DocumentReceivingMethod string `json:"document_receiving_method"`
			ReminderDelay           int    `json:"reminder_delay"`
			ReminderReceivingMethod string `json:"reminder_receiving_method"`
			Address                 struct {
				Street            string `json:"street"`
				AddressComplement string `json:"address_complement"`
				City              string `json:"city"`
				PostalCode        string `json:"postal_code"`
				Country           string `json:"country"`
			} `json:"address"`
		} `json:"occupants"`
		Units []struct {
			CadastralReference string `json:"cadastralReference"`
			Status             string `json:"status"`
			UnitType           string `json:"unitType"`
			Floor              uint8  `json:"floor"`
			Description        string `json:"description"`
			UnitAddress        struct {
				Street     string `json:"street"`
				Complement string `json:"complement"`
				City       string `json:"city"`
				PostalCode string `json:"postal_code"`
				Country    string `json:"country"`
			} `json:"unitAddress"`
			Owners []struct {
				Name          string `json:"name"`
				Surname       string `json:"surname"`
				Title         string `json:"title"`
				Quota         int    `json:"quota"`
				Administrator bool   `json:"administrator"`
			} `json:"owners"`
		} `json:"units"`
	}

	// Validation des données reçues
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data" + err.Error()})
		return
	}

	db := config.DB

	// Transaction pour garantir la cohérence des données
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ajouter l'adresse principale
	address := models.Address{
		Street:     requestData.Address.Street,
		Complement: requestData.Address.AddressComplement,
		City:       requestData.Address.City,
		PostalCode: requestData.Address.PostalCode,
		Country:    requestData.Address.Country,
	}

	if err := tx.Create(&address).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}
	var conciergeID *uint

	// Ajouter les occupants
	for _, occupantData := range requestData.Occupants {
		// Ajouter l'adresse du domicile de l'occupant
		domicileAddress := models.Address{
			Street:     occupantData.Address.Street,
			Complement: occupantData.Address.AddressComplement,
			City:       occupantData.Address.City,
			PostalCode: occupantData.Address.PostalCode,
			Country:    occupantData.Address.Country,
		}

		if err := tx.Create(&domicileAddress).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create occupant address" + err.Error()})
			return
		}

		// Ajouter l'occupant
		occupant := models.Occupant{
			Name:                      occupantData.Name,
			Surname:                   occupantData.Surname,
			Email:                     occupantData.Email,
			Corporation:               occupantData.Corporation,
			Phone:                     occupantData.Phone,
			Iban:                      occupantData.Iban,
			BirthDate:                 time.Now(),
			CivilityID:                1,
			DomicileAddressID:         domicileAddress.ID,
			DocumentReceivingMethodID: 1,
			ReminderDelay:             occupantData.ReminderDelay,
			ReminderReceivingMethodID: 1,
		}

		if err := tx.Create(&occupant).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create occupant" + err.Error()})
			return
		}
		if occupantData.IsConcierge {
			conciergeID = &occupant.ID
		}
	}

	// Ajouter le condominium
	condominium := models.Condominium{
		Name:        requestData.Informations.Name,
		Prefix:      requestData.Informations.Prefix,
		Description: requestData.Informations.Description,
		AddressID:   address.ID,
		ConciergeID: conciergeID,
	}

	if err := tx.Create(&condominium).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create condominium: " + err.Error() + " Address ID: " + fmt.Sprint(address.ID)})
		return
	}

	// Ajouter les unités
	for _, unitData := range requestData.Units {
		// Ajouter l'adresse de l'unité
		unitAddress := models.Address{
			Street:     unitData.UnitAddress.Street,
			Complement: unitData.UnitAddress.Complement,
			City:       unitData.UnitAddress.City,
			PostalCode: unitData.UnitAddress.PostalCode,
			Country:    unitData.UnitAddress.Country,
		}

		if err := tx.Create(&unitAddress).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create unit address"})
			return
		}

		var unitType models.UnitType
		if err := tx.Where("label = ?", unitData.UnitType).First(&unitType).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find unit type: " + err.Error()})
			return
		}

		// Ajouter l'unité
		unit := models.Unit{
			CondominiumID:      condominium.ID,
			AddressID:          unitAddress.ID,
			CadastralReference: unitData.CadastralReference,
			UnitTypeID:         unitType.ID,
			Floor:              unitData.Floor,
		}

		if err := tx.Create(&unit).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create unit" + err.Error()})
			return
		}

		// Ajouter les propriétaires de l'unité
		for _, ownerData := range unitData.Owners {
			label, err := getOccupantTypeLabel(ownerData.Title)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid occupant type: " + err.Error()})
				return
			}

			var occupantType models.OccupantType
			if err := tx.Where("label = ?", label).First(&occupantType).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find unit type: " + err.Error()})
				return
			}
			var occupant models.Occupant
			if err := tx.Where("name = ? AND surname = ?", ownerData.Name, ownerData.Surname).First(&occupant).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find occupant: " + err.Error()})
				return
			}

			owner := models.OccupantPossessionOnUnit{
				OccupantID:     occupant.ID,
				UnitID:         unit.ID,
				Quota:          float64(ownerData.Quota),
				Administrator:  ownerData.Administrator,
				OccupantTypeID: occupantType.ID,
			}

			if err := tx.Create(&owner).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create unit owner" + err.Error()})
				return
			}
		}
	}

	// Commit de la transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Condominium created successfully"})
}

func CheckUniqueness(c *gin.Context) {
	db := config.DB

	// Récupérer les paramètres de la requête
	name := c.Query("name")
	prefix := c.Query("prefix")

	var exists bool
	var err error

	// Vérifier l'unicité du nom
	if name != "" {
		conditions := map[string]interface{}{
			"name": name,
		}
		exists, err = CheckIfExists(db, "condominia", conditions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking name uniqueness"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"isTaken": exists})
		return
	}

	// Vérifier l'unicité du préfixe
	if prefix != "" {
		conditions := map[string]interface{}{
			"prefix": prefix,
		}
		exists, err = CheckIfExists(db, "condominia", conditions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking prefix uniqueness"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"isTaken": exists})
		return
	}

	// Si ni `name` ni `prefix` n'ont été fournis
	c.JSON(http.StatusBadRequest, gin.H{"error": "Either name or prefix must be provided"})
}
