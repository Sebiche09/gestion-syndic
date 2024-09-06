package upload

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Sebiche09/gestion-syndic/controller/errors"

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
		// Utilisation correcte de errors.handleError au lieu de HandleErrorc
		errors.HandleError(c, nil, ErrFileTypeMissing, http.StatusBadRequest)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		errors.HandleError(c, fmt.Errorf("erreur lors de la récupération du fichier : %w", err), ErrFileMissing, http.StatusInternalServerError)
		return
	}

	if !isValidFileType(fileType, file) {
		errors.HandleError(c, nil, ErrInvalidFileType, http.StatusBadRequest)
		return
	}

	switch fileType {
	case FileTypeCadastre:
		err = handleCadastreUpload(c, file)
		if err != nil {
			errors.HandleError(c, fmt.Errorf("échec du traitement du fichier cadastre : %w", err), "Erreur lors du traitement du fichier", http.StatusInternalServerError)
			return
		}
	default:
		errors.HandleError(c, nil, ErrUnknownFileType, http.StatusBadRequest)
	}
}

// HandleCadastreUpload gère l'upload de fichiers cadastre
func handleCadastreUpload(c *gin.Context, file *multipart.FileHeader) error {
	// Envoyer le fichier à PaddleOCR
	ocrResult, err := sendToPaddleOCR(file)
	if err != nil {
		// Ajouter du contexte à l'erreur renvoyée
		return fmt.Errorf("échec lors de l'envoi du fichier à PaddleOCR : %w", err)
	}

	cadastralData := extractCadastralData(ocrResult["text"].(string))
	// Retourner le résultat OCR
	c.JSON(http.StatusOK, gin.H{
		"message": "Fichier cadastre uploadé avec succès",
		"text":    cadastralData,
	})

	return nil
}

// isValidFileType vérifie si le fichier est du bon type selon son extension
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
