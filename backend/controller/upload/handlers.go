package upload

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Sebiche09/gestion-syndic/controller/errors"

	"github.com/gin-gonic/gin"
)

const (
	FileTypeCadastre   = "cadastre"                 // Type de fichier cadastre accepté
	FileTypePlan       = "plan"                     // Type de fichier plan accepté
	ErrFileTypeMissing = "Type de fichier manquant" // Erreur si le type de fichier n'est pas fourni
	ErrFileMissing     = "Fichier manquant"         // Erreur si aucun fichier n'est fourni
	ErrInvalidFileType = "Fichier de type invalide" // Erreur si le fichier n'a pas une extension valide
	ErrUnknownFileType = "Type de fichier inconnu"  // Erreur si le type de fichier est inconnu
)

// UploadHandler gère la logique principale pour l'upload de fichiers.
// Il vérifie d'abord que le type et le fichier sont bien fournis, puis valide le type du fichier.
// Si les validations réussissent, le fichier est traité en fonction de son type.
func UploadHandler(c *gin.Context) {
	log.Print("UploadHandler")
	fileType := c.PostForm("type")
	log.Print(fileType)
	if fileType == "" {
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
	case FileTypePlan:
		err = handlePlanUpload(c, file)
		if err != nil {
			errors.HandleError(c, fmt.Errorf("échec du traitement du fichier plan : %w", err), "Erreur lors du traitement du fichier", http.StatusInternalServerError)
			return
		}
	default:
		errors.HandleError(c, nil, ErrUnknownFileType, http.StatusBadRequest)
	}
}

// handleCadastreUpload gère l'upload et le traitement des fichiers cadastre.
// Il envoie le fichier à un service OCR (PaddleOCR), extrait les données cadastrales et retourne les résultats.
func handleCadastreUpload(c *gin.Context, file *multipart.FileHeader) error {
	tempFilePath := "/tmp/" + file.Filename
	err := c.SaveUploadedFile(file, tempFilePath)
	if err != nil {
	}
	ocrResult, err := sendToPaddleOCR(file)
	if err != nil {
		return fmt.Errorf("échec lors de l'envoi du fichier à PaddleOCR : %w", err)
	}

	cadastralData := extractCadastralData(ocrResult["text"].(string))

	c.JSON(http.StatusOK, gin.H{
		"message":  "Fichier cadastre uploadé avec succès",
		"text":     cadastralData,
		"filePath": tempFilePath,
	})

	return nil
}

// handleCadastreUpload gère l'upload et le traitement des fichiers cadastre.
// Il envoie le fichier à un service OCR (PaddleOCR), extrait les données cadastrales et retourne les résultats.
func handlePlanUpload(c *gin.Context, file *multipart.FileHeader) error {
	tempFilePath := "/tmp/" + file.Filename
	err := c.SaveUploadedFile(file, tempFilePath)
	if err != nil {
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Fichier cadastre uploadé avec succès",
		"filePath": tempFilePath,
	})

	return nil
}

// isValidFileType vérifie si le fichier uploadé a une extension valide en fonction de son type.
// Il accepte les fichiers PDF pour les types "pdf" et "cadastre".
func isValidFileType(fileType string, fileHeader *multipart.FileHeader) bool {
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))

	switch fileType {
	case FileTypeCadastre:
		return fileExt == ".pdf"
	case FileTypePlan:
		return fileExt == ".pdf"
	default:
		return false
	}
}
