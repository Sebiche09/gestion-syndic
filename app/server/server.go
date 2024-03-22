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
	d, err := sql.Open("mysql", dataSource()) //Ouvre une connexion à la DB.
	//sql.Open() renvoie objet de type *sql.DB, fct dataSource() founrit la chaine de connexion à la db
	if err != nil {
		log.Fatal(err)
	} // Si une erreur -> l'app s'arrete grace = log.Fatal
	defer d.Close() //defer permet de fermer la connexion avec la DB lorsque la fonction main est terminée

	// Créer une instance de moteur Gin
	r := gin.Default() //Crée une instance du moteur Gin (framework de golang)

	// Définir les routes
	setupRoutes(r, db.NewDB(d)) //Appelle la fct setupRoutes() avec en param l'instance gin et l'instance de la DB

	// CORS est activé uniquement en mode de production
	cors := os.Getenv("profile") == "prod"

	// Utiliser le moteur Gin
	addr := ":8080" // Port d'écoute
	if os.Getenv("profile") == "prod" {
		addr = ":80" // Port standard pour les applications web si on est en mode production (port 80)
	}
	log.Printf("Serveur démarré sur %s\n", addr) //print que le serveur est démarré
	if err := r.Run(addr); err != nil {          //Démarre le serveur et retourne log si erreur
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
} //dataSource retourne chaine de connexion formatée poour mysql avec les valeurs host et pass

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
