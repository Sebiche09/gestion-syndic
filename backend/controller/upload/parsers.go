package upload

import (
	"regexp"
	"strings"
)

// extractCadastralData extrait les informations cadastrales du texte OCR
func extractCadastralData(ocrText string) map[string][]OwnerInfo {
	extractedData := make(map[string][]OwnerInfo)

	normalizedText := normalizeText(ocrText)

	natureDetailRegex := regexp.MustCompile(`ENTITÉ PRIV.#(.*?)(?:RÉSULTAT\s*:|\d+\s+INFORMATION\s+CADASTRALE\s+ET\s+PATRIMONIALE\s+DE\s+LA\s+PARCELLE)`)
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

// parseAddress parse l'adresse complète
func parseAddress(fullAddress string) AddressInfo {
	var address AddressInfo

	parts := strings.Split(fullAddress, " - ")
	if len(parts) == 2 {
		address.Street = strings.TrimSpace(parts[0])
		postalCityMatch := regexp.MustCompile(`(\d{4,5})\s+(.+)`).FindStringSubmatch(parts[1])
		if len(postalCityMatch) == 3 {
			address.PostalCode = postalCityMatch[1]
			address.City = postalCityMatch[2]
		} else {
			address.Country = strings.TrimSpace(parts[1])
		}
	} else {
		postalCityMatch := regexp.MustCompile(`(.*)\s+(\d{4,5})\s+(.+)`).FindStringSubmatch(fullAddress)
		if len(postalCityMatch) == 4 {
			address.Street = strings.TrimSpace(postalCityMatch[1])
			address.PostalCode = postalCityMatch[2]
			address.City = postalCityMatch[3]
		} else {
			address.Street = fullAddress
		}
	}

	return address
}

// normalizeText normalise le texte OCR
func normalizeText(text string) string {
	normalized := strings.ReplaceAll(text, "\n", " ")
	normalized = strings.ReplaceAll(normalized, "\r", "")
	return strings.Join(strings.Fields(normalized), " ")
}
