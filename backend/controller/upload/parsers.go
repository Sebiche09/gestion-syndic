package upload

import (
	"regexp"
	"strings"
)

// extractMainAddress extrait l'adresse principale du bâtiment à partir du texte OCR.
// Si aucun pays n'est précisé, on ajoute "Belgique" à l'adresse.
func extractMainAddress(ocrText string) AddressInfo {
	// Expression régulière pour capturer l'adresse entre les deux balises.
	addressRegex := regexp.MustCompile(`INFORMATION\s+CADASTRALE\s+ET\s+PATRIMONIALE\s+DE\s+LA\s+PARCELLE\s+([\s\S]+?)\s+Section\s+et\s+n°\s+de\s+parcelle`)
	addressMatch := addressRegex.FindStringSubmatch(ocrText)

	// Si une correspondance est trouvée, on procède à l'extraction des détails de l'adresse.
	if len(addressMatch) > 1 {
		fullAddress := strings.TrimSpace(addressMatch[1])

		// On capture deux parties : la rue avec le numéro et la ville avec le code postal.
		addressParts := regexp.MustCompile(`(.+?)\s+(\d{4,5})\s+([A-Za-zÀ-ÖØ-öø-ÿ -]+)`).FindStringSubmatch(fullAddress)

		if len(addressParts) > 3 {
			// Crée un objet AddressInfo avec les informations extraites.
			address := AddressInfo{
				Street:     strings.TrimSpace(addressParts[1]),                // Rue + numéro
				PostalCode: strings.TrimSpace(addressParts[2]),                // Code postal
				City:       cleanCityName(strings.TrimSpace(addressParts[3])), // Ville
			}

			// Si aucun pays n'est précisé, on ajoute "Belgique".
			// On peut vérifier s'il y a un pays dans la ville. Si non, on ajoute Belgique.
			if !strings.Contains(ocrText, "FR") {
				address.Country = "BE"
			}

			// Retourne l'adresse principale extraite.
			return address
		}
	}

	// Si l'adresse n'est pas trouvée, on retourne une adresse vide.
	return AddressInfo{}
}

// extractCadastralData extrait les informations cadastrales d'un texte OCR.
// Il s'appuie sur des expressions régulières pour capturer les identifiants de parcelles
// et associer les propriétaires trouvés dans le texte au bon identifiant.
func extractCadastralData(ocrText string) map[string]interface{} {
	// Crée une map pour stocker les informations extraites.
	extractedData := make(map[string]interface{})

	// Normalise le texte pour supprimer les retours à la ligne et les espaces multiples.
	normalizedText := normalizeText(ocrText)

	// Extrait l'adresse principale du texte OCR.
	mainAddress := extractMainAddress(normalizedText)
	extractedData["address"] = mainAddress

	// Expression régulière pour capturer les identifiants après le # dans le texte OCR.
	natureDetailRegex := regexp.MustCompile(`Fin\s+exonération[\s\S]+?#([\s\S]+?)(?:RÉSULTAT\s*:|\d{1,2}\s+INFORMATION\s+CADASTRALE\s+ET\s+PATRIMONIALE\s+DE\s+LA\s+PARCELLE)`)
	matches := natureDetailRegex.FindAllStringSubmatch(normalizedText, -1)

	// Map pour stocker les informations des parcelles et propriétaires.
	unit := make(map[string][]OwnerInfo)

	// Boucle sur chaque correspondance pour extraire les informations.
	for _, match := range matches {
		if len(match) > 1 {
			fullDetail := strings.TrimSpace(match[1])
			lines := strings.Split(fullDetail, " ")

			if len(lines) > 0 {
				identifier := strings.ReplaceAll(strings.TrimSpace(lines[0]), " ", "")

				if strings.HasPrefix(identifier, "Cave") {
					caveKey := strings.TrimSpace(lines[1])
					fullKey := "Cave " + caveKey
					owners := extractOwners(fullDetail)
					unit[fullKey] = owners
				} else {
					owners := extractOwners(fullDetail)
					unit[identifier] = owners
				}
			}
		}
	}

	// Ajoute les informations des parcelles au résultat final.
	extractedData["unit"] = unit

	// Retourne le résultat final avec l'adresse principale et les informations cadastrales.
	return extractedData
}

// extractOwners extrait et renvoie les informations des propriétaires à partir des détails complets capturés.
// Utilise une expression régulière pour capturer les noms, adresses et titres des propriétaires.
func extractOwners(fullDetail string) []OwnerInfo {
	var owners []OwnerInfo

	// Expression régulière pour capturer les informations des propriétaires.
	// Elle récupère le nom, l'adresse, le code postal, la ville, et d'autres titres (PP, NP, etc.).
	ownerRegex := regexp.MustCompile(`(\d+)\s+([A-Za-zÀ-ÖØ-öø-ÿ' -]+),\s*([A-Za-zÀ-ÖØ-öø-ÿ' -]+)\s+((?:Rue|Avenue|Boulevard|Chemin|Place|Chaussée|Route|Clos).+?|-\s+)(?:\s+-\s+(\d{4,5}))?\s*([A-Za-zÀ-ÖØ-öø-ÿ' -]+)?\s*(PP\s+\d+/\d+|NP\s+\d+/\d+|US\s+\d+/\d+|Ust\s+\d+/\d+)`)
	ownerMatches := ownerRegex.FindAllStringSubmatch(fullDetail, -1)

	// Boucle sur chaque correspondance pour construire les objets OwnerInfo.
	for _, match := range ownerMatches {
		if len(match) >= 7 {
			var street string
			// Si la rue est manquante, remplace par un message.
			if strings.TrimSpace(match[4]) == "-" {
				street = "Adresse non disponible"
			} else {
				street = strings.TrimSpace(match[4])
			}

			// Crée l'adresse avec les informations extraites.
			address := AddressInfo{
				Street:     street,
				PostalCode: strings.TrimSpace(match[5]),
				City:       cleanCityName(strings.TrimSpace(match[6])),
			}

			// Vérifie si le pays est vide, si c'est le cas, assigne "Belgique" par défaut.
			if address.Country == "" {
				address.Country = "Belgique"
			}

			// Crée l'objet OwnerInfo avec les données extraites.
			owner := OwnerInfo{
				LastName:  strings.TrimSpace(match[2]),
				FirstName: strings.TrimSpace(match[3]),
				Address:   address,
				Title:     strings.TrimSpace(match[7]),
			}

			// Ajoute le propriétaire à la liste.
			owners = append(owners, owner)
		}
	}

	// Retourne la liste des propriétaires.
	return owners
}

// parseAddress analyse une adresse complète et en extrait les différentes parties (rue, code postal, ville).
func parseAddress(fullAddress string) AddressInfo {
	var address AddressInfo

	// Expression régulière pour capturer les différentes parties d'une adresse : rue, code postal, ville, etc.
	postalCityMatch := regexp.MustCompile(`(.*?)(?:\s*\((.*?)\))?\s+(\d{3,5})(?:\s+Bte\s+\w+)?\s+-\s+(\d{4,5})\s+([A-Za-zÀ-ÖØ-öø-ÿ -]+)`).FindStringSubmatch(fullAddress)
	if len(postalCityMatch) == 6 {
		street := postalCityMatch[1]
		if postalCityMatch[2] != "" {
			street += " (" + postalCityMatch[2] + ")"
		}
		street += " " + postalCityMatch[3]

		// Si une boîte postale ("Bte") est présente, on l'ajoute à la rue.
		if strings.Contains(fullAddress, "Bte") {
			bteMatch := regexp.MustCompile(`Bte\s+\w+`).FindString(fullAddress)
			street += " " + bteMatch
		}

		// Assigne les parties de l'adresse.
		address.Street = strings.TrimSpace(street)
		address.PostalCode = postalCityMatch[4]
		address.City = cleanCityName(postalCityMatch[5])
	} else {
		// Si la regex échoue, tout attribuer à la rue.
		address.Street = cleanExtraInfo(fullAddress)
	}

	// Retourne l'objet AddressInfo avec les informations analysées.
	return address
}

// cleanCityName nettoie les informations superflues qui suivent le nom de la ville.
func cleanCityName(city string) string {
	// Liste de mots-clés qui indiquent la fin du nom de la ville et le début des informations superflues.
	invalidInfo := []string{"PP", "NP", "US", "Ust", "SUPERF", "USA/HAB", "EMPH", "PROPRIÉTAIRE", "PROPRIETAIRE", "FRANCE"}

	// Boucle pour retirer ces informations si elles sont présentes dans la ville.
	for _, word := range invalidInfo {
		if idx := strings.Index(city, word); idx != -1 {
			return strings.TrimSpace(city[:idx])
		}
	}

	// Retourne le nom de la ville nettoyé.
	return city
}

// cleanExtraInfo supprime les informations non pertinentes à la fin des adresses.
func cleanExtraInfo(street string) string {
	// Liste de mots-clés qui indiquent la fin des informations pertinentes dans une adresse.
	invalidInfo := []string{"PP", "NP", "US", "PROPRIÉTAIRE", "PROPRIETAIRE"}

	// Boucle pour enlever les parties superflues si elles sont présentes dans l'adresse.
	for _, word := range invalidInfo {
		if idx := strings.Index(street, word); idx != -1 {
			return strings.TrimSpace(street[:idx])
		}
	}

	// Retourne l'adresse nettoyée.
	return street
}

// normalizeText nettoie le texte OCR en supprimant les retours à la ligne et les espaces multiples.
func normalizeText(text string) string {
	// Remplace les retours à la ligne par des espaces.
	normalized := strings.ReplaceAll(text, "\n", " ")
	normalized = strings.ReplaceAll(normalized, "\r", "")

	// Supprime les espaces multiples pour uniformiser le texte.
	return strings.Join(strings.Fields(normalized), " ")
}
