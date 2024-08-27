package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Sebiche09/gestion-syndic/config"
	"github.com/Sebiche09/gestion-syndic/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getIdByType(db *gorm.DB, tableName, typeName string) (int, error) {
	var result struct {
		ID int
	}

	err := db.Table(tableName).Select("id").Where("type = ?", typeName).Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.ID, nil
}

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

func CreateCondominium(c *gin.Context) {
	var requestData struct {
		Informations struct {
			Name        string `json:"name"`
			Prefix      string `json:"prefix"`
			Description string `json:"description"`
		} `json:"informations"`
		Address struct {
			Street            string `json:"street"`
			Number            string `json:"number"`
			AddressComplement string `json:"address_complement"`
			City              string `json:"city"`
			PostalCode        string `json:"postal_code"`
			Country           string `json:"country"`
		} `json:"address"`
		FtpBlueprint struct {
			Blueprint string `json:"blueprint"`
		} `json:"ftpBlueprint"`
		Concierge struct {
			Name                       string `json:"name"`
			Surname                    string `json:"surname"`
			Email                      string `json:"email"`
			Corporation                bool   `json:"corporation"`
			Phone                      string `json:"phone"`
			Iban                       string `json:"iban"`
			BirthDate                  string `json:"birthdate"` // Consider using time.Time and parsing it
			Civility                   string `json:"civility"`
			DocumentReceivingMethod    string `json:"document_receiving_method"`
			ReminderDelay              string `json:"reminder_delay"`
			ReminderReceivingMethod    string `json:"reminder_receiving_method"`
			Street                     string `json:"street_concierge"`
			Number                     string `json:"number_concierge"`
			AddressComplementConcierge string `json:"address_complement_concierge"`
			City                       string `json:"city_concierge"`
			PostalCode                 string `json:"postal_code_concierge"`
			Country                    string `json:"country_concierge"`
		} `json:"concierge"`
	}

	db := config.DB

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminderDelay, err := strconv.Atoi(requestData.Concierge.ReminderDelay)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder_delay format"})
		return
	}
	birthDate, err := time.Parse("2006-01-02", requestData.Concierge.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birthdate format"})
		return
	}

	// Vérifier si un condominium existe déjà avec le même nom
	conditionsCondominiumName := map[string]interface{}{
		"name": requestData.Informations.Name,
	}
	existsNameCondominium, err := CheckIfExists(db, "condominia", conditionsCondominiumName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking condominium name existence"})
		return
	}
	if existsNameCondominium {
		c.JSON(http.StatusConflict, gin.H{"error": "Condominium with this name already exists"})
		return
	}

	// Vérifier si un condominium existe déjà avec le même préfixe
	conditionsCondominiumPrefix := map[string]interface{}{
		"prefix": requestData.Informations.Prefix,
	}
	existsPrefixCondominium, err := CheckIfExists(db, "condominia", conditionsCondominiumPrefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking condominium prefix existence"})
		return
	}
	if existsPrefixCondominium {
		c.JSON(http.StatusConflict, gin.H{"error": "Condominium with this prefix already exists"})
		return
	}

	// Vérifier si un occupant existe deja avec le même nom et prénom et date de naissance
	conditionsOccupant := map[string]interface{}{
		"name":       requestData.Concierge.Name,
		"surname":    requestData.Concierge.Surname,
		"birth_date": birthDate,
	}
	existsOccupant, err := CheckIfExists(db, "occupants", conditionsOccupant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error checking occupant existence"})
		return
	}
	if existsOccupant {
		c.JSON(http.StatusConflict, gin.H{"Error": "this occupant already exists"})
		return
	}

	// Insert Address
	address := models.Address{
		Street:     requestData.Address.Street,
		Number:     requestData.Address.Number,
		Complement: requestData.Address.AddressComplement,
		City:       requestData.Address.City,
		PostalCode: requestData.Address.PostalCode,
		Country:    requestData.Address.Country,
	}
	if err := db.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}
	address_concierge := models.Address{
		Street:     requestData.Concierge.Street,
		Number:     requestData.Concierge.Number,
		Complement: requestData.Concierge.AddressComplementConcierge,
		City:       requestData.Concierge.City,
		PostalCode: requestData.Concierge.PostalCode,
		Country:    requestData.Concierge.Country,
	}
	if err := db.Create(&address_concierge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create concierge address"})
		return
	}

	civilityID, err := getIdByType(db, "civilities", requestData.Concierge.Civility)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find civility"})
		return
	}
	documentReceivingMethodID, err := getIdByType(db, "document_receiving_methods", requestData.Concierge.DocumentReceivingMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find document receiving method"})
		return
	}
	reminderReceivingMethodID, err := getIdByType(db, "reminder_receiving_methods", requestData.Concierge.ReminderReceivingMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find reminder receiving method"})
		return
	}

	occupant := models.Occupant{
		Name:                      requestData.Concierge.Name,
		Surname:                   requestData.Concierge.Surname,
		Email:                     requestData.Concierge.Email,
		Corporation:               requestData.Concierge.Corporation,
		Phone:                     requestData.Concierge.Phone,
		Iban:                      requestData.Concierge.Iban,
		BirthDate:                 birthDate,
		CivilityID:                civilityID,
		DomicileAddressID:         address_concierge.ID,
		DocumentReceivingMethodID: documentReceivingMethodID,
		ReminderDelay:             reminderDelay,
		ReminderReceivingMethodID: reminderReceivingMethodID,
	}

	if err := db.Create(&occupant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create occupant"})
		return
	}

	// Insert Condominium
	condominium := models.Condominium{
		Name:               requestData.Informations.Name,
		AddressID:          address.ID,
		Description:        requestData.Informations.Description,
		FtpBlueprintPath:   requestData.FtpBlueprint.Blueprint,
		LandRegistryNumber: requestData.Informations.Prefix,
		Prefix:             requestData.Informations.Prefix,
		ConciergeID:        occupant.ID,
	}

	if err := db.Create(&condominium).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create condominium"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Condominium created successfully"})
}

func GetCivilities(c *gin.Context) {
	var civilities []models.Civility
	if err := config.DB.Find(&civilities).Error; err != nil {
		handleError(c, err, "Error fetching civilities", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, civilities)
}
func GetDocumentReceivingMethods(c *gin.Context) {
	var documentReceivingMethods []models.DocumentReceivingMethod
	if err := config.DB.Find(&documentReceivingMethods).Error; err != nil {
		handleError(c, err, "Error fetching document receiving methods", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, documentReceivingMethods)
}
func GetReminderReceivingMethods(c *gin.Context) {
	var reminderReceivingMethods []models.ReminderReceivingMethod
	if err := config.DB.Find(&reminderReceivingMethods).Error; err != nil {
		handleError(c, err, "Error fetching reminder receiving methods", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, reminderReceivingMethods)
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
