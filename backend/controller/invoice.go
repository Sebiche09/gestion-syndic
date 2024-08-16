package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Sebiche09/gestion-syndic/models"
	"github.com/gin-gonic/gin"
	"github.com/jlaffaye/ftp"
	"gorm.io/gorm"
)

// UploadHandler gère les requêtes d'upload de fichier
func UploadHandler(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		handleError(c, err, "Error parsing form data", http.StatusBadRequest)
		return
	}

	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		handleError(c, err, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Envoyer le fichier à PaddleOCR
	ocrResult, err := sendToPaddleOCR(file, handler.Filename)
	if err != nil {
		handleError(c, err, "Error sending request to OCR service", http.StatusInternalServerError)
		return
	}

	// Extraire le texte de la réponse OCR
	text, ok := ocrResult["text"].(string)
	if !ok {
		handleError(c, nil, "Error extracting text from OCR response", http.StatusInternalServerError)
		return
	}
	prompt := `Je vais te fournir un texte extrait d'une facture. À partir de ce texte, je souhaite que tu extrayes les informations pertinentes et que tu les formates en un objet JSON correspondant à la structure suivante:
	go
	type Invoice struct {
		gorm.Model
		InvoiceType       bool 
		InvoiceNumber     string      
		InternalReference string      
		InvoiceLabel      string      
		InvoiceDate       time.Time   
		SupplierID        uint        
		Supplier          Supplier  
		CondominiumID     uint       
		Condominium       Condominium 
		PriceInclVAT      float64  
		PriceExclVAT      float64 
		InvoiceStatus     uint     
		FtpFilePath       string
		ContractID        uint
		Contract          Contract 
		ExerciceID        uint 
		Exercice          Exercice 
	} voici le texte : ` + text
	model := "llama3"
	// Envoyer le texte extrait à Ollama
	ollamaResult, err := sendToOllama(prompt, model)
	if err != nil {
		handleError(c, err, "Error sending request to Mistral", http.StatusInternalServerError)
		return
	}

	// Envoyer les données extraites à la base de données
	c.JSON(http.StatusOK, gin.H{"data": ollamaResult})
}

// Fonction pour envoyer le fichier à PaddleOCR
func sendToPaddleOCR(file io.Reader, filename string) (map[string]interface{}, error) {
	//Permet de structurer la requête pour envoyer le fichier à PaddleOCR
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	writer.Close()
	// Envoi de la requête à PaddleOCR
	resp, err := http.Post("http://ocr:5000/ocr", writer.FormDataContentType(), &requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Décoder la réponse de PaddleOCR
	var ocrResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&ocrResult)
	if err != nil {
		return nil, err
	}

	return ocrResult, nil
}

// Fonction pour envoyer le texte extrait à Ollama
func sendToOllama(prompt string, model string) (string, error) {
	requestData := map[string]string{"prompt": prompt, "model": model}
	// Convertir les données en JSON
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request data: %w", err)
	}
	// Envoi de la requête à Ollama
	resp, err := http.Post("http://ollama:11434/api/generate", "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var responseMap map[string]interface{}
	fullResponse := ""
	// Boucle pour récupérer la réponse partielle
	for {
		// Decode la réponse
		if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
			return "", fmt.Errorf("failed to decode response: %w", err)
		}

		// Collecter la réponse partielle
		if responsePart, ok := responseMap["response"].(string); ok {
			fullResponse += responsePart
		}

		// Vérifier si la génération est terminée
		if done, ok := responseMap["done"].(bool); ok && done {
			break
		}
	}

	return fullResponse, nil
}

// UploadHandler gère les requêtes d'upload de fichier
func sendToFTP(c *gin.Context) {
	// Parse the multipart form with a maximum upload size of 10 MB
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form data"})
		log.Printf("Error parsing form data: %v", err) // Log the error
		return
	}

	// Retrieve the file from the form data
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving file"})
		log.Printf("Error retrieving file: %v", err) // Log the error
		return
	}
	defer file.Close()
	// Read FTP configuration from environment variables
	ftpServer := os.Getenv("FTP_SERVER")
	ftpUser := os.Getenv("FTP_USER")
	ftpPassword := os.Getenv("FTP_PASSWORD")
	ftpDir := "/temporary/pending"

	// Connect to the FTP server
	ftpConn, err := ftp.Dial(ftpServer, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to FTP server"})
		log.Printf("Error connecting to FTP server: %v", err) // Log the error
		return
	}
	defer ftpConn.Quit()

	// Authenticate with the FTP server
	err = ftpConn.Login(ftpUser, ftpPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging into FTP server"})
		log.Printf("Error logging into FTP server: %v", err) // Log the error
		return
	}

	// Change to the desired directory on the FTP server
	err = ftpConn.ChangeDir(ftpDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error changing FTP directory"})
		log.Printf("Error changing FTP directory: %v", err) // Log the error
		return
	}

	// Upload the file to the FTP server
	ftpFilePath := filepath.Join(ftpDir, handler.Filename)
	err = ftpConn.Stor(ftpFilePath, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading file to FTP server"})
		log.Printf("Error uploading file to FTP server: %v", err) // Log the error
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File uploaded successfully to FTP: %s", handler.Filename)})
}

// Fonction pour envoyer le texte extrait à la base de données
func sendToDatabase(text string, c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Déclaration de la variable pour stocker les données extraites
	var invoiceData map[string]interface{}

	// Décoder le texte JSON pour obtenir les données
	err := json.Unmarshal([]byte(text), &invoiceData)
	if err != nil {
		handleError(c, err, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	// Convertir les données en types appropriés
	invoiceDateStr, _ := invoiceData["InvoiceDate"].(string)
	invoiceDate, err := time.Parse(time.RFC3339, invoiceDateStr)
	if err != nil {
		handleError(c, err, "Error parsing date", http.StatusInternalServerError)
		return
	}

	supplierData, _ := invoiceData["Supplier"].(map[string]interface{})
	var supplier models.Supplier
	if name, ok := supplierData["Name"].(string); ok {
		supplier.Name = name
	}

	// Créer ou mettre à jour le fournisseur
	result := db

	// Créer l'objet Invoice
	invoice := models.Invoice{
		InvoiceType:       invoiceData["InvoiceType"].(bool),
		InvoiceNumber:     invoiceData["InvoiceNumber"].(string),
		InternalReference: invoiceData["InternalReference"].(string),
		InvoiceLabel:      invoiceData["InvoiceLabel"].(string),
		InvoiceDate:       invoiceDate,
		SupplierID:        supplier.ID,
		CondominiumID:     0, // Not provided in JSON, set to default
		PriceInclVAT:      invoiceData["PriceInclVAT"].(float64),
		PriceExclVAT:      invoiceData["PriceExclVAT"].(float64),
		InvoiceStatus:     uint(invoiceData["InvoiceStatus"].(float64)), // Assuming it's a number
		FtpFilePath:       invoiceData["FtpFilePath"].(string),
		ContractID:        0, // Not provided in JSON, set to default
		ExerciceID:        0, // Not provided in JSON, set to default
	}

	// Enregistrer l'objet Invoice dans la base de données
	result = db.Create(&invoice)
	if result.Error != nil {
		handleError(c, result.Error, "Error saving invoice to database", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data successfully saved to the database"})
}

// Fonction pour gérer les erreurs
func handleError(c *gin.Context, err error, message string, statusCode int) {
	log.Printf("Error: %v", err)
	c.JSON(statusCode, gin.H{"error": message})
}
