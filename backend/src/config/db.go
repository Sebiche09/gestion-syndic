package config

import (
	"log"
	"time"

	"github.com/Sebiche09/gestion-syndic/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SeedDatabase() {
	var count int64

	// Insert default civilities
	if err := DB.Model(&models.Civility{}).Count(&count).Error; err != nil {
		log.Fatal("Failed to count civility entries: ", err)
	}
	if count == 0 {
		civilities := []models.Civility{
			{Type: "Monsieur"},
			{Type: "Madame"},
		}
		if err := DB.Create(&civilities).Error; err != nil {
			log.Fatal("Failed to seed civility data: ", err)
		}
		log.Println("Civility data seeded successfully")
	}

	// Insert default document receiving methods
	if err := DB.Model(&models.DocumentReceivingMethod{}).Count(&count).Error; err != nil {
		log.Fatal("Failed to count document receiving method entries: ", err)
	}
	if count == 0 {
		documentReceivingMethods := []models.DocumentReceivingMethod{
			{Type: "Email"},
			{Type: "Courrier"},
			{Type: "Fax"},
			{Type: "Recommandé"},
		}
		if err := DB.Create(&documentReceivingMethods).Error; err != nil {
			log.Fatal("Failed to seed receiving method data: ", err)
		}
		log.Println("Document receiving method data seeded successfully")
	}

	// Insert default reminder receiving methods
	if err := DB.Model(&models.ReminderReceivingMethod{}).Count(&count).Error; err != nil {
		log.Fatal("Failed to count reminder receiving method entries: ", err)
	}
	if count == 0 {
		reminderReceivingMethods := []models.ReminderReceivingMethod{
			{Type: "Email"},
			{Type: "Courrier"},
			{Type: "SMS"},
			{Type: "Fax"},
			{Type: "Recommandé"},
		}
		if err := DB.Create(&reminderReceivingMethods).Error; err != nil {
			log.Fatal("Failed to seed reminder receiving method data: ", err)
		}
		log.Println("Reminder receiving method data seeded successfully")
	}

	// Insert default unit types
	if err := DB.Model(&models.UnitType{}).Count(&count).Error; err != nil {
		log.Fatal("Failed to count unit type entries: ", err)
	}
	if count == 0 {
		unitTypes := []models.UnitType{
			{Label: "Appartement"},
			{Label: "Cave"},
			{Label: "Garage"},
			{Label: "Local Commercial"},
			{Label: "Autre"},
		}
		if err := DB.Create(&unitTypes).Error; err != nil {
			log.Fatal("Failed to seed unit type data: ", err)
		}
		log.Println("Unit type data seeded successfully")
	}
	if err := DB.Model(&models.OccupantType{}).Count(&count).Error; err != nil {
		log.Fatal("Failed to count occupant type entries: ", err)
	}
	if count == 0 {
		occupantTypes := []models.OccupantType{
			{Label: "pleine propriete"},
			{Label: "nue propriete"},
			{Label: "usufruit"},
			{Label: "superficiaire"},
			{Label: "emphyteote"},
			{Label: "usage/habitation"},
		}
		if err := DB.Create(&occupantTypes).Error; err != nil {
			log.Fatal("Failed to seed occupant type data: ", err)
		}
		log.Println("Occupant type data seeded successfully")
	}
}

func Connect() {
	dsn := "user=postgres dbname=postgres password=postgres host=db port=5432"
	var db *gorm.DB
	var err error

	for i := 1; i <= 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Successfully connected to the database!")
			DB = db
			break
		}

		log.Printf("Failed to connect to database (attempt %d/10): %v", i, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after 10 attempts: ", err)
	}

	err = db.AutoMigrate(
		&models.Occupant{},
		&models.Civility{},
		&models.Address{},
		&models.DocumentReceivingMethod{},
		&models.ReminderReceivingMethod{},
		&models.OccupantPossessionOnUnit{},
		&models.OccupantType{},
		&models.Unit{},
		&models.ElectricGazMeter{},
		&models.UnitType{},
		&models.Condominium{},
		&models.Exercice{},
		&models.Contract{},
		&models.BankAccount{},
		&models.BankAccountType{},
		&models.AccountStatement{},
		&models.AccountStatementOccupant{},
		&models.GeneralAssembly{},
		&models.GAParticipation{},
		&models.CondoSupplier{},
		&models.Supplier{},
		&models.SupplierCategory{},
		&models.Invoice{},
		&models.AccountStatementInvoice{},
		&models.Remender{},
		&models.AllocationKey{},
		&models.AllocationKeyTemplate{},
		&models.KeyUnit{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}

	log.Println("Database migration completed successfully")
	SeedDatabase()
}
