package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Sebiche09/gestion-syndic/config"
	"github.com/Sebiche09/gestion-syndic/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	router := gin.New()
	config.Connect()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.OccupantRoute(router)
	routes.InvoiceRoute(router)
	routes.CondominiumRoute(router)
	routes.CivilityRoute(router)
	routes.ReceivingMethodRoute(router)
	router.Run(":8080")
}

func loadEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	envFile := fmt.Sprintf(".env.%s", env)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}
}
