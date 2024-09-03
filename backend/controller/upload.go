package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// Struct pour stocker les informations d'un propriétaire
type OwnerInfo struct {
	LastName  string
	FirstName string
	Address   AddressInfo
	Title     string
}

// Struct pour stocker les informations d'adresse
type AddressInfo struct {
	Street     string
	PostalCode string
	City       string
	Country    string
}

// Gérer les téléchargements de fichiers
func UploadHandler(c *gin.Context) {
	// Récupérer le type de fichier à partir du formulaire
	fileType := c.PostForm("type")
	if fileType == "" {
		handleError(c, nil, "Type de fichier manquant", http.StatusBadRequest)
		return
	}

	// Récupérer le fichier à partir du formulaire
	file, err := c.FormFile("file")
	if err != nil {
		handleError(c, err, "Fichier manquant", http.StatusInternalServerError)
		return
	}

	// Vérifier que le fichier est un PDF si c'est le type attendu
	if fileType == "pdf" && !isValidPDF(file) {
		handleError(c, nil, "Fichier PDF invalide", http.StatusBadRequest)
		return
	}

	// Gérer le fichier selon son type
	switch fileType {
	case "cadastre":
		handleCadastreUpload(c, file)
	}

}

// Fonction pour vérifier si le fichier est un PDF
func isValidPDF(fileHeader *multipart.FileHeader) bool {
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	return fileExt == ".pdf"
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

// Fonction pour extraire les informations cadastrales
func extractCadastralData(ocrText string) map[string][]OwnerInfo {
	extractedData := make(map[string][]OwnerInfo)

	// Normaliser le texte pour enlever les espaces et sauts de ligne superflus
	normalizedText := strings.ReplaceAll(ocrText, "\n", " ")
	normalizedText = strings.ReplaceAll(normalizedText, "\r", "")
	normalizedText = strings.Join(strings.Fields(normalizedText), " ")

	// Expression régulière pour capturer tout entre "ENTITÉ PRIV.#" et "RÉSULTAT :" ou "n INFORMATION CADASTRALE ET PATRIMONIALE DE LA PARCELLE"
	natureDetailRegex := regexp.MustCompile(`ENTITÉ PRIV.#(.*?)(?:RÉSULTAT\s*:|\d+\s+INFORMATION\s+CADASTRALE\s+ET\s+PATRIMONIALE\s+DE\s+LA\s+PARCELLE)`)
	matches := natureDetailRegex.FindAllStringSubmatch(normalizedText, -1)

	for _, match := range matches {
		if len(match) > 1 {
			// Capturer tout le texte entre "ENTITÉ PRIV.#" et "RÉSULTAT :" ou "n INFORMATION CADASTRALE ET PATRIMONIALE DE LA PARCELLE"
			fullDetail := strings.TrimSpace(match[1])

			// Extraire une clé unique pour cette entrée, par exemple, en utilisant la première ligne
			lines := strings.Split(fullDetail, " ")
			if len(lines) > 0 {
				natureDetail := strings.TrimSpace(lines[0])
				owners := extractOwners(fullDetail)
				extractedData[natureDetail] = owners
			}
		}
	}

	return extractedData
}

// Fonction pour extraire les informations des propriétaires d'un lot
func extractOwners(fullDetail string) []OwnerInfo {
	var owners []OwnerInfo

	// Expression régulière pour trouver chaque propriétaire
	ownerRegex := regexp.MustCompile(`(\d+)\s+(\w+),\s+(\w+)\s+(.*?)\s+(PP\s+\d+/\d+|NP\s+\d+/\d+|US\s+\d+/\d+)`)
	ownerMatches := ownerRegex.FindAllStringSubmatch(fullDetail, -1)

	for _, match := range ownerMatches {
		if len(match) > 5 {
			address := parseAddress(match[4])
			owner := OwnerInfo{
				LastName:  strings.TrimSpace(match[2]),
				FirstName: strings.TrimSpace(match[3]),
				Address:   address,
				Title:     strings.TrimSpace(match[5]),
			}
			owners = append(owners, owner)
		}
	}

	return owners
}

// Fonction pour parser l'adresse en composants : rue, code postal, ville, pays
func parseAddress(fullAddress string) AddressInfo {
	// Initialiser les composants de l'adresse
	var address AddressInfo

	// Vérifier si l'adresse contient un tiret séparant la rue et le reste
	parts := strings.Split(fullAddress, " - ")
	if len(parts) == 2 {
		address.Street = strings.TrimSpace(parts[0])
		// Vérifier si la partie droite contient un code postal suivi d'une ville
		postalCityMatch := regexp.MustCompile(`(\d{4,5})\s+(.+)`).FindStringSubmatch(parts[1])
		if len(postalCityMatch) == 3 {
			address.PostalCode = postalCityMatch[1]
			address.City = postalCityMatch[2]
		} else {
			address.Country = strings.TrimSpace(parts[1]) // Sinon, c'est un pays
		}
	} else {
		// Pas de tiret, on essaye de capturer une adresse classique avec code postal et ville
		postalCityMatch := regexp.MustCompile(`(.*)\s+(\d{4,5})\s+(.+)`).FindStringSubmatch(fullAddress)
		if len(postalCityMatch) == 4 {
			address.Street = strings.TrimSpace(postalCityMatch[1])
			address.PostalCode = postalCityMatch[2]
			address.City = postalCityMatch[3]
		} else {
			// Sinon, tout est dans la rue
			address.Street = fullAddress
		}
	}

	return address
}

// Envoyer le fichier à PaddleOCR pour l'OCR
func sendToPaddleOCR(fileHeader *multipart.FileHeader) (map[string]interface{}, error) {
	// Ouvrir le fichier pour le lire
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Structurer la requête pour envoyer le fichier à PaddleOCR
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	// Copier le contenu du fichier dans le writer
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	writer.Close()

	// Envoyer la requête à PaddleOCR avec le bon type MIME (multipart/form-data)
	resp, err := http.Post("http://ocr:5000/ocr", writer.FormDataContentType(), &requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Lire le corps de la réponse pour capturer le message d'erreur
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	// Décoder la réponse de PaddleOCR
	var ocrResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&ocrResult)
	if err != nil {
		return nil, err
	}
	return ocrResult, nil
}

// Gérer les erreurs
func handleError(c *gin.Context, err error, message string, statusCode int) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
	c.JSON(statusCode, gin.H{"error": message})
}
