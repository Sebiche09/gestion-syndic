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
	FileTypePDF        = "pdf"                      // Type de fichier PDF accepté
	FileTypeCadastre   = "cadastre"                 // Type de fichier cadastre accepté
	ErrFileTypeMissing = "Type de fichier manquant" // Erreur si le type de fichier n'est pas fourni
	ErrFileMissing     = "Fichier manquant"         // Erreur si aucun fichier n'est fourni
	ErrInvalidFileType = "Fichier de type invalide" // Erreur si le fichier n'a pas une extension valide
	ErrUnknownFileType = "Type de fichier inconnu"  // Erreur si le type de fichier est inconnu
)

// UploadHandler gère la logique principale pour l'upload de fichiers.
// Il vérifie d'abord que le type et le fichier sont bien fournis, puis valide le type du fichier.
// Si les validations réussissent, le fichier est traité en fonction de son type.
func UploadHandler(c *gin.Context) {
	// Récupère le type de fichier depuis le formulaire
	fileType := c.PostForm("type")
	if fileType == "" {
		// Si aucun type de fichier n'est fourni, une erreur est retournée
		errors.HandleError(c, nil, ErrFileTypeMissing, http.StatusBadRequest)
		return
	}

	// Récupère le fichier à partir du formulaire
	file, err := c.FormFile("file")
	if err != nil {
		// Si le fichier ne peut pas être récupéré, une erreur serveur est retournée
		errors.HandleError(c, fmt.Errorf("erreur lors de la récupération du fichier : %w", err), ErrFileMissing, http.StatusInternalServerError)
		return
	}

	// Vérifie si le fichier a une extension valide pour le type de fichier donné
	if !isValidFileType(fileType, file) {
		// Si le type de fichier est invalide, une erreur est retournée
		errors.HandleError(c, nil, ErrInvalidFileType, http.StatusBadRequest)
		return
	}

	// Traite le fichier selon son type
	switch fileType {
	case FileTypeCadastre:
		// Si c'est un fichier cadastre, il est envoyé pour un traitement OCR
		err = handleCadastreUpload(c, file)
		if err != nil {
			// Si le traitement échoue, une erreur est retournée
			errors.HandleError(c, fmt.Errorf("échec du traitement du fichier cadastre : %w", err), "Erreur lors du traitement du fichier", http.StatusInternalServerError)
			return
		}
	default:
		// Si le type de fichier est inconnu, une erreur est retournée
		errors.HandleError(c, nil, ErrUnknownFileType, http.StatusBadRequest)
	}
}

// handleCadastreUpload gère l'upload et le traitement des fichiers cadastre.
// Il envoie le fichier à un service OCR (PaddleOCR), extrait les données cadastrales et retourne les résultats.
func handleCadastreUpload(c *gin.Context, file *multipart.FileHeader) error {
	// Envoie le fichier au service PaddleOCR pour analyse OCR
	ocrResult, err := sendToPaddleOCR(file)
	if err != nil {
		// Si l'envoi échoue, retourne une erreur avec un message descriptif
		return fmt.Errorf("échec lors de l'envoi du fichier à PaddleOCR : %w", err)
	}

	// Extrait les informations cadastrales à partir du résultat OCR
	cadastralData := extractCadastralData(ocrResult["text"].(string))

	// Retourne les données OCR au client sous forme de JSON
	c.JSON(http.StatusOK, gin.H{
		"message": "Fichier cadastre uploadé avec succès",
		"text":    cadastralData,
	})

	// Si tout s'est bien passé, retourne nil
	return nil
}

// isValidFileType vérifie si le fichier uploadé a une extension valide en fonction de son type.
// Il accepte les fichiers PDF pour les types "pdf" et "cadastre".
func isValidFileType(fileType string, fileHeader *multipart.FileHeader) bool {
	// Extrait l'extension du fichier (en minuscule pour standardiser)
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))

	// Vérifie l'extension en fonction du type de fichier
	switch fileType {
	case FileTypeCadastre:
		// Si le type est cadastre, l'extension doit aussi être .pdf
		return fileExt == ".pdf"
	default:
		// Si le type de fichier est inconnu, retourne false
		return false
	}
}
