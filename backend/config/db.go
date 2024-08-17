package config

import (
	"log"
	"time"

	"github.com/Sebiche09/gestion-syndic/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SeedDatabase() {
	// Vérifier si les données de base existent déjà pour éviter les duplications
	var count int64

	// Insert default civility types
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

	// Insert default receiving methods
	if err := DB.Model(&models.ReceivingMethod{}).Count(&count).Error; err != nil {
		log.Fatal("Failed to count receiving method entries: ", err)
	}
	if count == 0 {
		receivingMethods := []models.ReceivingMethod{
			{Type: "Email"},
			{Type: "Courrier"},
			{Type: "SMS"},
			{Type: "Fax"},
			{Type: "Recommandé"},
		}
		if err := DB.Create(&receivingMethods).Error; err != nil {
			log.Fatal("Failed to seed receiving method data: ", err)
		}
		log.Println("Receiving method data seeded successfully")
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
		&models.ReceivingMethod{},
		&models.OccupantPossessionOnProperty{},
		&models.OccupantType{},
		&models.Property{},
		&models.ElectricGazMeter{},
		&models.PropertyType{},
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
		&models.KeyProperty{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}

	log.Println("Database migration completed successfully")
	SeedDatabase()
}
