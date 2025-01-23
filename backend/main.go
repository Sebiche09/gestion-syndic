package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/Sebiche09/gestion-syndic/docs"
	"github.com/Sebiche09/gestion-syndic/src/config"
	"github.com/Sebiche09/gestion-syndic/src/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	loadEnv()
	router := gin.New()
	config.Connect()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	router.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))
	routes.OccupantRoute(router)
	routes.InvoiceRoute(router)
	routes.CondominiumRoute(router)
	routes.UnitRoute(router)
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
