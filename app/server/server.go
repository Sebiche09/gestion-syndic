package main

import (
	"database/sql"
	"github.com/gin-gonic/gin" // Importation de Gin
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"sogim/db"
)

func main() {
	d, err := sql.Open("mysql", dataSource())
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	// Créer une instance de moteur Gin
	r := gin.Default()

	// Définir les routes
	setupRoutes(r, db.NewDB(d))

	// CORS est activé uniquement en mode de production
	cors := os.Getenv("profile") == "prod"

	// Utiliser le moteur Gin
	addr := ":8080" // Port d'écoute
	if os.Getenv("profile") == "prod" {
		addr = ":80" // Port standard pour les applications web
	}
	log.Printf("Serveur démarré sur %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func dataSource() string {
	host := "localhost"
	pass := "pass"
	if os.Getenv("profile") == "prod" {
		host = "db"
		pass = os.Getenv("db_pass")
	}
	return "goxygen:" + pass + "@tcp(" + host + ":3306)/goxygen"
}

// setupRoutes définit les routes de l'application web avec Gin
func setupRoutes(r *gin.Engine, db *db.DB) {
	// Définir une route GET
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// Autres routes...
}
