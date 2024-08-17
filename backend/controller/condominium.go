package controller

import (
	"net/http"
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
			ReminderDelay              int    `json:"reminder_delay"`
			ReminderReceivingMethod    string `json:"reminder_receiving_method"`
			Street                     string `json:"street_concierge"`
			Number                     string `json:"number_concierge"`
			AddressComplementConcierge string `json:"address_complement_concierge"`
			City                       string `json:"city_concierge"`
			PostalCode                 string `json:"postal_code_concierge"`
			Country                    string `json:"country_concierge"`
		} `json:"concierge"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

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
	reminderReceivingMethodID, err := getIdByType(db, "reminder_receiving_method", requestData.Concierge.ReminderReceivingMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find reminder receiving method"})
		return
	}
	birthDate, err := time.Parse("2006-01-02", requestData.Concierge.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birthdate format"})
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
		ReminderDelay:             requestData.Concierge.ReminderDelay,
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
