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

	"github.com/gin-gonic/gin"
	"github.com/jlaffaye/ftp"
)

// Fonction pour gérer les erreurs
func handleError(c *gin.Context, err error, message string, statusCode int) {
	log.Printf("Error: %v", err)
	c.JSON(statusCode, gin.H{"error": message})
}

// Fonction pour envoyer le fichier à PaddleOCR
func sendToPaddleOCR(file io.Reader, filename string) (map[string]interface{}, error) {
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

	resp, err := http.Post("http://ocr:5000/ocr", writer.FormDataContentType(), &requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request data: %w", err)
	}

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
	for {
		if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
			return "", fmt.Errorf("failed to decode response: %w", err)
		}

		// Collect the response part
		if responsePart, ok := responseMap["response"].(string); ok {
			fullResponse += responsePart
		}

		// Check if the response is complete
		if done, ok := responseMap["done"].(bool); ok && done {
			break
		}
		// Optionally wait a bit before polling again (implement a delay or use another approach if needed)
	}

	return fullResponse, nil
}

// Handler pour le endpoint de téléchargement
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
	prompt := "peux tu récupérer le numéro de facture, le nom du fournisseur et me le structurer en tiret à partir du texte suivant : " + text
	model := "llama3"
	// Envoyer le texte extrait à Ollama
	ollamaResult, err := sendToOllama(prompt, model)
	if err != nil {
		handleError(c, err, "Error sending request to Mistral", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ollamaResult)
}

// UploadHandler gère les requêtes d'upload de fichier
func UploadHandlers(c *gin.Context) {
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
