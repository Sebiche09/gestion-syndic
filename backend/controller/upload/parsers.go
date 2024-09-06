package upload

import (
	"log"
	"regexp"
	"strings"
)

// extractCadastralData extrait les informations cadastrales du texte OCR
func extractCadastralData(ocrText string) map[string][]OwnerInfo {
	extractedData := make(map[string][]OwnerInfo)

	normalizedText := normalizeText(ocrText)

	// Mise à jour de la regex pour capturer un ou deux chiffres avant 'INFORMATION CADASTRALE...'
	natureDetailRegex := regexp.MustCompile(`ENTITÉ PRIV.#(.*?)(?:RÉSULTAT\s*:|\d{1,2}\s+INFORMATION\s+CADASTRALE\s+ET\s+PATRIMONIALE\s+DE\s+LA\s+PARCELLE)`)
	matches := natureDetailRegex.FindAllStringSubmatch(normalizedText, -1)

	for _, match := range matches {
		if len(match) > 1 {
			fullDetail := strings.TrimSpace(match[1])
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

// extractOwners extrait les propriétaires du détail complet
func extractOwners(fullDetail string) []OwnerInfo {
	var owners []OwnerInfo

	// Regex amélioré pour capturer les informations du propriétaire avec titre
	ownerRegex := regexp.MustCompile(`(\d+)\s+([A-Za-zÀ-ÖØ-öø-ÿ' -]+),\s*([A-Za-zÀ-ÖØ-öø-ÿ' -]+)\s+((?:Rue|Avenue|Boulevard|Chemin|Place|Chaussée|R|RUE|ROUTE|Route).+?)\s+-\s+(\d{4,5})?\s*([A-Za-zÀ-ÖØ-öø-ÿ' -]+)?\s*(PP\s+\d+/\d+|NP\s+\d+/\d+|US\s+\d+/\d+)`)
	log.Print(fullDetail)
	ownerMatches := ownerRegex.FindAllStringSubmatch(fullDetail, -1)

	for _, match := range ownerMatches {
		// Vérification des groupes capturés pour assurer que le titre est présent
		if len(match) >= 7 {
			address := AddressInfo{
				Street:     strings.TrimSpace(match[4]),
				PostalCode: strings.TrimSpace(match[5]),
				City:       cleanCityName(strings.TrimSpace(match[6])),
			}
			owner := OwnerInfo{
				LastName:  strings.TrimSpace(match[2]),
				FirstName: strings.TrimSpace(match[3]),
				Address:   address,
				Title:     strings.TrimSpace(match[7]), // Capturer le titre ici
			}
			owners = append(owners, owner)
		}
	}

	return owners
}

// parseAddress parse l'adresse complète
func parseAddress(fullAddress string) AddressInfo {
	var address AddressInfo

	// Regex mise à jour pour capturer une rue, numéro, Bte, code postal, et ville
	postalCityMatch := regexp.MustCompile(`(.*?)(?:\s*\((.*?)\))?\s+(\d{3,5})(?:\s+Bte\s+\w+)?\s+-\s+(\d{4,5})\s+([A-Za-zÀ-ÖØ-öø-ÿ -]+)`).FindStringSubmatch(fullAddress)
	if len(postalCityMatch) == 6 {
		street := postalCityMatch[1] // Partie de l'adresse avant les parenthèses
		if postalCityMatch[2] != "" {
			street += " (" + postalCityMatch[2] + ")" // Ajout des parenthèses si elles existent
		}
		street += " " + postalCityMatch[3] // Ajout du numéro de rue

		// Si la boîte postale ("Bte") est présente, on la garde dans la rue
		if strings.Contains(fullAddress, "Bte") {
			bteMatch := regexp.MustCompile(`Bte\s+\w+`).FindString(fullAddress)
			street += " " + bteMatch
		}

		// On assigne les parties de l'adresse
		address.Street = strings.TrimSpace(street)
		address.PostalCode = postalCityMatch[4]
		address.City = cleanCityName(postalCityMatch[5])
	} else {
		// Si la regex ne fonctionne pas, on attribue tout à la rue
		address.Street = cleanExtraInfo(fullAddress)
	}

	return address
}

// cleanCityName nettoie les informations superflues après la ville
func cleanCityName(city string) string {
	// Liste des mots-clés qui peuvent indiquer la fin de la ville et l'apparition d'informations superflues
	invalidInfo := []string{"PP", "NP", "US", "PROPRIÉTAIRE", "PROPRIETAIRE", "FRANCE"}

	for _, word := range invalidInfo {
		if idx := strings.Index(city, word); idx != -1 {
			return strings.TrimSpace(city[:idx])
		}
	}

	return city
}

// cleanExtraInfo supprime les informations non liées à l'adresse (comme PP 1/1)
func cleanExtraInfo(street string) string {
	// Liste des mots-clés qui peuvent indiquer la fin de l'adresse
	invalidInfo := []string{"PP", "NP", "US", "PROPRIÉTAIRE", "PROPRIETAIRE"}

	for _, word := range invalidInfo {
		if idx := strings.Index(street, word); idx != -1 {
			return strings.TrimSpace(street[:idx])
		}
	}

	return street
}

// normalizeText normalise le texte OCR en supprimant les retours à la ligne inutiles et les espaces multiples
func normalizeText(text string) string {
	// Remplacer les sauts de ligne par des espaces
	normalized := strings.ReplaceAll(text, "\n", " ")
	normalized = strings.ReplaceAll(normalized, "\r", "")

	// Supprimer les espaces multiples
	return strings.Join(strings.Fields(normalized), " ")
}
