package upload

import (
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	FileTypePDF        = "pdf"
	FileTypeCadastre   = "cadastre"
	ErrFileTypeMissing = "Type de fichier manquant"
	ErrFileMissing     = "Fichier manquant"
	ErrInvalidFileType = "Fichier de type invalide"
	ErrUnknownFileType = "Type de fichier inconnu"
)

// UploadHandler gère les uploads de fichiers
func UploadHandler(c *gin.Context) {
	fileType := c.PostForm("type")
	if fileType == "" {
		handleError(c, nil, ErrFileTypeMissing, http.StatusBadRequest)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		handleError(c, err, ErrFileMissing, http.StatusInternalServerError)
		return
	}

	if !isValidFileType(fileType, file) {
		handleError(c, nil, ErrInvalidFileType, http.StatusBadRequest)
		return
	}

	switch fileType {
	case FileTypeCadastre:
		handleCadastreUpload(c, file)
	default:
		handleError(c, nil, ErrUnknownFileType, http.StatusBadRequest)
	}
	log.Print("UploadHandler terminé")
}

// isValidFileType vérifie si le fichier est du bon type basé sur son extension
func isValidFileType(fileType string, fileHeader *multipart.FileHeader) bool {
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))

	switch fileType {
	case FileTypePDF:
		return fileExt == ".pdf"
	case FileTypeCadastre:
		return fileExt == ".pdf" // Ou une autre extension si nécessaire
	default:
		return false
	}
}

// Gestion spécifique pour les fichiers PDF
func handleCadastreUpload(c *gin.Context, file *multipart.FileHeader) {
	ocrResult, err := sendToPaddleOCR(file)
	if err != nil {
		handleError(c, err, "Erreur lors de l'envoi à PaddleOCR", http.StatusInternalServerError)
		return
	}

	// Extraire les informations spécifiques du résultat OCR
	cadastralData := extractCadastralData(ocrResult["text"].(string))
	// Retourner le résultat de l'OCR
	c.JSON(http.StatusOK, gin.H{
		"message": "PDF upload handled",
		"text":    cadastralData,
	})
}
