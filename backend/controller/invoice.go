package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jlaffaye/ftp"
)

// UploadHandler gère les requêtes d'upload de fichier
func UploadHandler(c *gin.Context) {
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
