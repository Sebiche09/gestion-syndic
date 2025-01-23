# Gestion des Condominiums

Ce fichier contient les fonctions liées à la gestion des condominiums, notamment la récupération, la création, et la vérification d'unicité. Voici une explication détaillée de chaque fonction.

---

## 1. `GetAllCondominiums`

### Description
Cette fonction récupère la liste de tous les condominiums enregistrés dans la base de données, ainsi que leurs adresses associées.

### Code :
```go
func GetAllCondominiums(c *gin.Context) {
    db := config.DB

    var condominiums []models.Condominium

    // Chargement de tous les condominiums avec leurs adresses associées
    if err := db.Preload("Address").Find(&condominiums).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch condominiums"})
        return
    }

    // Construire la réponse avec les champs nécessaires
    response := []gin.H{}
    for _, condo := range condominiums {
        response = append(response, gin.H{
            "name":   condo.Name,
            "prefix": condo.Prefix,
            "city":   condo.Address.City,
        })
    }
    c.JSON(http.StatusOK, response)
}
```
---

### Explications
```go
db := config.DB
```
*   Récupère l'instance de la base de données configurée dans config.DB.
*   Cela permet d'exécuter des opérations sur la base de données.

```go
var condominiums []models.Condominium
```
*   Déclare une variable condominiums comme un tableau de structures models.Condominium
*   Cette variable stockera les résultats de la requête à la base de données.

```go
if err := db.Preload("Address").Find(&condominiums).Error; err != nil{
```
*   Effectue une requête pour récupérer tous les condominiums depuis la base de données.
*   La méthode Preload("Address") charge également les adresses associées à chaque condominium (relation hasOne ou hasMany dans GORM)
*   Si une erreur survient pendant la requête, elle est stockée dans la variable err.

```go
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch condominiums"})
        return
```
*   Si err n’est pas nil (c’est-à-dire si une erreur s’est produite), cette ligne envoie une réponse HTTP avec un code d’erreur 500 Internal Server Error.
*   La réponse contient un message JSON précisant que la récupération des condominiums a échoué.
*   La fonction retourne immédiatement, arrêtant son exécution.

```go
response := []gin.H{}
```
*   Initialise une liste vide de type []gin.H qui servira à construire la réponse JSON.
*   Chaque élément de cette liste représentera un condominium sous forme de map clé-valeur.

```go
for _, condo := range condominiums {
```
*   Parcourt chaque élément de la liste condominiums retournée par la base de données.
*   condo représente un condominium individuel pour chaque itération.

```go
response = append(response, gin.H{
            "name":   condo.Name,
            "prefix": condo.Prefix,
            "city":   condo.Address.City,
        })
```
*   Ajoute un nouvel élément à la liste response.
*   Cet élément est une map (gin.H) contenant :
    *   name : Le nom du condominium (condo.Name).
    *   prefix : Le préfixe unique du condominium (condo.Prefix).
    *   city : La ville de l'adresse associée au condominium (condo.Address.City).

```go
c.JSON(http.StatusOK, response)
```
*   Envoie une réponse HTTP avec un code 200 OK.
*   Le corps de la réponse contient la liste JSON response, représentant tous les condominiums avec leurs informations pertinentes.

## 2. `CreateCondominium`

### Description
Cette fonction crée un nouveau condominium avec ses adresses, occupants et unités associés

### Code :
```go
func CreateCondominium(c *gin.Context) {
	var requestData struct {
		Informations struct {
			Name        string `json:"name"`
			Prefix      string `json:"prefix"`
			Description string `json:"description"`
		} `json:"informations"`
		Address struct {
			Street            string `json:"street"`
			AddressComplement string `json:"address_complement"`
			City              string `json:"city"`
			PostalCode        string `json:"postal_code"`
			Country           string `json:"country"`
		} `json:"address"`
		Occupants []struct {
			Name                    string `json:"name"`
			Surname                 string `json:"surname"`
			Email                   string `json:"email"`
			Corporation             bool   `json:"corporation"`
			IsConcierge             bool   `json:"isConcierge"`
			Phone                   string `json:"phone"`
			Iban                    string `json:"iban"`
			BirthDate               string `json:"birthdate"`
			Civility                string `json:"civility"`
			DocumentReceivingMethod string `json:"document_receiving_method"`
			ReminderDelay           int    `json:"reminder_delay"`
			ReminderReceivingMethod string `json:"reminder_receiving_method"`
			Address                 struct {
				Street            string `json:"street"`
				AddressComplement string `json:"address_complement"`
				City              string `json:"city"`
				PostalCode        string `json:"postal_code"`
				Country           string `json:"country"`
			} `json:"address"`
		} `json:"occupants"`
		Units []struct {
			CadastralReference string `json:"cadastralReference"`
			Status             string `json:"status"`
			UnitType           string `json:"unitType"`
			Floor              uint8  `json:"floor"`
			Description        string `json:"description"`
			UnitAddress        struct {
				Street     string `json:"street"`
				Complement string `json:"complement"`
				City       string `json:"city"`
				PostalCode string `json:"postal_code"`
				Country    string `json:"country"`
			} `json:"unitAddress"`
			Owners []struct {
				Name          string `json:"name"`
				Surname       string `json:"surname"`
				Title         string `json:"title"`
				Quota         int    `json:"quota"`
				Administrator bool   `json:"administrator"`
			} `json:"owners"`
		} `json:"units"`
	}

	// Validation des données reçues
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data" + err.Error()})
		return
	}

	db := config.DB

	// Transaction pour garantir la cohérence des données
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ajouter l'adresse principale
	address := models.Address{
		Street:     requestData.Address.Street,
		Complement: requestData.Address.AddressComplement,
		City:       requestData.Address.City,
		PostalCode: requestData.Address.PostalCode,
		Country:    requestData.Address.Country,
	}

	if err := tx.Create(&address).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}
	var conciergeID *uint

	// Ajouter les occupants
	for _, occupantData := range requestData.Occupants {
		// Ajouter l'adresse du domicile de l'occupant
		domicileAddress := models.Address{
			Street:     occupantData.Address.Street,
			Complement: occupantData.Address.AddressComplement,
			City:       occupantData.Address.City,
			PostalCode: occupantData.Address.PostalCode,
			Country:    occupantData.Address.Country,
		}

		if err := tx.Create(&domicileAddress).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create occupant address" + err.Error()})
			return
		}

		// Ajouter l'occupant
		occupant := models.Occupant{
			Name:                      occupantData.Name,
			Surname:                   occupantData.Surname,
			Email:                     occupantData.Email,
			Corporation:               occupantData.Corporation,
			Phone:                     occupantData.Phone,
			Iban:                      occupantData.Iban,
			BirthDate:                 time.Now(),
			CivilityID:                1,
			DomicileAddressID:         domicileAddress.ID,
			DocumentReceivingMethodID: 1,
			ReminderDelay:             occupantData.ReminderDelay,
			ReminderReceivingMethodID: 1,
		}

		if err := tx.Create(&occupant).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create occupant" + err.Error()})
			return
		}
		if occupantData.IsConcierge {
			conciergeID = &occupant.ID
		}
	}

	// Ajouter le condominium
	condominium := models.Condominium{
		Name:        requestData.Informations.Name,
		Prefix:      requestData.Informations.Prefix,
		Description: requestData.Informations.Description,
		AddressID:   address.ID,
		ConciergeID: conciergeID,
	}

	if err := tx.Create(&condominium).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create condominium: " + err.Error() + " Address ID: " + fmt.Sprint(address.ID)})
		return
	}

	// Ajouter les unités
	for _, unitData := range requestData.Units {
		// Ajouter l'adresse de l'unité
		unitAddress := models.Address{
			Street:     unitData.UnitAddress.Street,
			Complement: unitData.UnitAddress.Complement,
			City:       unitData.UnitAddress.City,
			PostalCode: unitData.UnitAddress.PostalCode,
			Country:    unitData.UnitAddress.Country,
		}

		if err := tx.Create(&unitAddress).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create unit address"})
			return
		}

		var unitType models.UnitType
		if err := tx.Where("label = ?", unitData.UnitType).First(&unitType).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find unit type: " + err.Error()})
			return
		}

		// Ajouter l'unité
		unit := models.Unit{
			CondominiumID:      condominium.ID,
			AddressID:          unitAddress.ID,
			CadastralReference: unitData.CadastralReference,
			UnitTypeID:         unitType.ID,
			Floor:              unitData.Floor,
		}

		if err := tx.Create(&unit).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create unit" + err.Error()})
			return
		}

		// Ajouter les propriétaires de l'unité
		for _, ownerData := range unitData.Owners {
			label, err := getOccupantTypeLabel(ownerData.Title)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid occupant type: " + err.Error()})
				return
			}

			var occupantType models.OccupantType
			if err := tx.Where("label = ?", label).First(&occupantType).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find unit type: " + err.Error()})
				return
			}
			var occupant models.Occupant
			if err := tx.Where("name = ? AND surname = ?", ownerData.Name, ownerData.Surname).First(&occupant).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find occupant: " + err.Error()})
				return
			}

			owner := models.OccupantPossessionOnUnit{
				OccupantID:     occupant.ID,
				UnitID:         unit.ID,
				Quota:          float64(ownerData.Quota),
				Administrator:  ownerData.Administrator,
				OccupantTypeID: occupantType.ID,
			}

			if err := tx.Create(&owner).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create unit owner" + err.Error()})
				return
			}
		}
	}

	// Commit de la transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Condominium created successfully"})
}
```


### Explications

#### 1. Définition de la structure de données pour la requête
```go
var requestData struct {
    ...
```
La structure requestData est utilisée pour mapper les données JSON de la requête HTTP. Cela inclut :
*   Informations : Contient les détails de base du condominium (nom, préfixe, description).
*   Address : L'adresse principale du condominium.
*   Occupants : Une liste des occupants avec leurs informations personnelles et leurs adresses.
*   Units : Une liste des unités avec leurs détails (adresse, propriétaires, etc.).

#### 2. Validation des données reçues
```go
if err := c.ShouldBindJSON(&requestData); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data" + err.Error()})
    return
}
```
*   La fonction ShouldBindJSON vérifie que les données reçues sont conformes à la structure définie dans requestData.
*   Si la validation échoue, une réponse HTTP avec le code 400 Bad Request est envoyée avec un message d’erreur.

#### 3. Initialisation de la transaction
```go
db := config.DB
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()
```
*   Récupère l'instance de la base de données via config.DB.
*   Démarre une transaction (tx) pour garantir que toutes les opérations liées à la création du condominium se déroulent de manière atomique.
*   Si une erreur survient ou si la fonction panique, la transaction est annulée avec Rollback.

#### 4. Création de l'adresse principale
```go
address := models.Address{
    Street:     requestData.Address.Street,
    Complement: requestData.Address.AddressComplement,
    City:       requestData.Address.City,
    PostalCode: requestData.Address.PostalCode,
    Country:    requestData.Address.Country,
}


if err := tx.Create(&address).Error; err != nil {
    tx.Rollback()
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
    return
}
```
*   Une structure Address est créée avec les données de l’adresse principale fournies dans la requête.
*   La méthode tx.Create insère cette adresse dans la base de données.
*   Si l’opération échoue, la transaction est annulée et une erreur HTTP 500 Internal Server Error est renvoyée.

#### 5. Ajout des occupants
```go
for _, occupantData := range requestData.Occupants {
    domicileAddress := models.Address{
        // Données de l'adresse du domicile
    }
    if err := tx.Create(&domicileAddress).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create occupant address"})
        return
    }

    occupant := models.Occupant{
        Name:          occupantData.Name,
        Surname:       occupantData.Surname,
        DomicileAddressID: domicileAddress.ID,
        // Autres champs
    }

    if err := tx.Create(&occupant).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create occupant"})
        return
    }
    if occupantData.IsConcierge {
        conciergeID = &occupant.ID
    }
}
```
*   Pour chaque occupant, une nouvelle adresse de domicile est créée et insérée dans la base de données.
*   Chaque occupant est ensuite créé avec une référence à l’adresse de son domicile.
*   Si un occupant est marqué comme concierge (IsConcierge), son ID est enregistré.

#### 6. Création du condominium
```go
condominium := models.Condominium{
    Name:        requestData.Informations.Name,
    Prefix:      requestData.Informations.Prefix,
    Description: requestData.Informations.Description,
    AddressID:   address.ID,
    ConciergeID: conciergeID,
}

if err := tx.Create(&condominium).Error; err != nil {
    tx.Rollback()
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create condominium"})
    return
}
```
*   Crée une nouvelle instance de Condominium avec les informations fournies.
*   L’adresse principale (AddressID) et l’ID du concierge (ConciergeID) sont associés.

#### 7.  Ajout des unités et propriétaires
```go
for _, unitData := range requestData.Units {
    for _, ownerData := range unitData.Owners {
    }
}
```
*   Chaque unité est insérée dans la base de données avec une adresse unique.
*   Les propriétaires de chaque unité sont également insérés avec leurs quotas et rôles.

#### 8. Validation et retour de la réponse
```go
if err := tx.Commit().Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
    return
}

c.JSON(http.StatusOK, gin.H{"message": "Condominium created successfully"})
```
*   Si toutes les étapes ont réussi, la transaction est validée avec Commit.
*   Une réponse HTTP avec le code 200 OK est renvoyée pour confirmer la création du condominium.

## 3. `CheckUniqueness`

### Description
Cette fonction vérifie si un nom ou un préfixe est déjà utilisé pour un condominium

### Code : 
```go
func CheckUniqueness(c *gin.Context) {
	db := config.DB

	// Récupérer les paramètres de la requête
	name := c.Query("name")
	prefix := c.Query("prefix")

	var exists bool
	var err error

	// Vérifier l'unicité du nom
	if name != "" {
		conditions := map[string]interface{}{
			"name": name,
		}
		exists, err = CheckIfExists(db, "condominia", conditions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking name uniqueness"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"isTaken": exists})
		return
	}

	// Vérifier l'unicité du préfixe
	if prefix != "" {
		conditions := map[string]interface{}{
			"prefix": prefix,
		}
		exists, err = CheckIfExists(db, "condominia", conditions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking prefix uniqueness"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"isTaken": exists})
		return
	}

	// Si ni `name` ni `prefix` n'ont été fournis
	c.JSON(http.StatusBadRequest, gin.H{"error": "Either name or prefix must be provided"})
}
```

### Explications

```go
db := config.DB
```
*   Récupère l'instance de la base de données configurée dans config.DB. Cette connexion est utilisée pour interroger la table condominia.

```go
name := c.Query("name")
prefix := c.Query("prefix")
```
*   Lit les paramètres de la requête HTTP (name et prefix), qui sont passés dans l'URL sous la forme de ?name=...&prefix=....
*   Ces paramètres sont utilisés pour vérifier l'unicité dans la base de données.
*   Cas possibles:
    *   Si name est fourni, la vérification s'effectuera sur le champ name.
    *   Si prefix est fourni, la vérification s'effectuera sur le champ prefix.

```go
var exists bool
var err error
```
*   exists : Indique si la valeur recherchée existe dans la base de données (true ou false).
*   err : Capture les éventuelles erreurs lors de l'interrogation de la base de données

```go
if name != "" {
    conditions := map[string]interface{}{
        "name": name,
    }
    exists, err = CheckIfExists(db, "condominia", conditions)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking name uniqueness"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"isTaken": exists})
    return
}
```
*   Si le paramètre name est fourni, on crée une map de conditions avec la clé name.
*   La fonction utilitaire CheckIfExists est appelée pour vérifier si une entrée correspondante existe dans la table condominia.
*   Gestion des erreurs :
    *   Si CheckIfExists retourne une erreur (err n’est pas nil), une réponse HTTP 500 Internal Server Error est envoyée avec un message d’erreur.µ
*   Réponse : 
    *   Si la vérification réussit, une réponse HTTP 200 OK est envoyée avec un champ isTaken indiquant si le nom est pris (true ou false).
*   Après avoir traité name, la fonction retourne immédiatement.

```go
if prefix != "" {
    conditions := map[string]interface{}{
        "prefix": prefix,
    }
    exists, err = CheckIfExists(db, "condominia", conditions)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking prefix uniqueness"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"isTaken": exists})
    return
}
```
*   Si le paramètre prefix est fourni, le processus est similaire à celui de name :
    *   Une map de conditions est créée avec la clé prefix.
    *   La fonction CheckIfExists vérifie si une entrée correspondante existe dans la table condominia.
*   Gestion des erreurs :
    *   Si une erreur survient, la réponse HTTP 500 est renvoyée.
*   Réponse :
    *   Retourne une réponse HTTP 200 OK avec un champ isTaken.

```go
c.JSON(http.StatusBadRequest, gin.H{"error": "Either name or prefix must be provided"})
```
*   Si ni name ni prefix ne sont fournis, une réponse HTTP 400 Bad Request est envoyée avec un message d’erreur.
*   Cela garantit que l’API ne tente pas de vérifier sans avoir les informations nécessaires.






## 4. `CheckIfExists`

### Description
Cette fonction vérifie si une valeur existe déjà dans une table donnée.

### Code : 
```go
func CheckIfExists(db *gorm.DB, tableName string, conditions map[string]interface{}) (bool, error) {
	var count int64
	query := db.Table(tableName).Where(conditions).Count(&count)
	if query.Error != nil {
		return false, query.Error
	}
	return count > 0, nil
}
```

### Explications

1.  Spécifie la table (tableName) et les conditions (conditions)
2.  Exécute une requête COUNT pour compter les entrées correspondant aux conditions.
3.  Retourne true si le compteur est supérieur à 0, sinon false.

## 5. `getOccupantTypeLabel`

### Description
Cette fonction retourne le label correspondant à une abréviation trouvée dans un titre.

### Code : 
```go
func getOccupantTypeLabel(title string) (string, error) {
	// Mapping des abréviations aux labels
	abbreviationToLabel := map[string]string{
		"PP":      "pleine propriete",
		"NP":      "nue propriete",
		"US":      "usufruit",
		"SUPERF":  "superficiaire",
		"USA/HAB": "usage/habitation",
		"emph":    "emphyteote",
	}

	// Vérification si le titre contient une des abréviations
	for abbreviation, label := range abbreviationToLabel {
		if strings.Contains(strings.ToUpper(title), abbreviation) {
			return label, nil
		}
	}

	// Retourner une erreur si aucune correspondance n'est trouvée
	return "", fmt.Errorf("no matching occupant type for title: %s", title)
}
```

### Explications
```go
abbreviationToLabel := map[string]string{
    "PP":      "pleine propriete",
    "NP":      "nue propriete",
    "US":      "usufruit",
    "SUPERF":  "superficiaire",
    "USA/HAB": "usage/habitation",
    "emph":    "emphyteote",
}
```
*   Une map (abbreviationToLabel) est utilisée pour associer des abréviations (clés) à leurs labels correspondants (valeurs).
*   Pourquoi une map ?
    *   Une map est idéale pour rechercher rapidement une valeur à partir d'une clé.
    *   Elle permet de réduire la complexité de recherche par rapport à une liste ou un tableau.

```go
for abbreviation, label := range abbreviationToLabel {
    if strings.Contains(strings.ToUpper(title), abbreviation) {
        return label, nil
    }
}
```
*   La boucle for parcourt chaque paire clé-valeur (abbreviation et label) de la map abbreviationToLabel.
*   La condition strings.Contains(strings.ToUpper(title), abbreviation) vérifie si l'abréviation actuelle est présente dans title :
    *   strings.ToUpper(title) : Convertit le titre en majuscules pour rendre la recherche insensible à la casse.
    *   strings.Contains : Vérifie si abbreviation est une sous-chaîne de title.
*   Si une correspondance est trouvée, la fonction retourne immédiatement :
    *   label : Le label correspondant à l'abréviation trouvée.
    *   nil : Aucune erreur.

```go
return "", fmt.Errorf("no matching occupant type for title: %s", title)
```
*   Si la boucle ne trouve aucune correspondance, la fonction retourne :
    *   Une chaîne vide ("") comme valeur pour le label.
    *   Une erreur formatée avec fmt.Errorf contenant un message indiquant qu'aucun type correspondant n'a été trouvé, ainsi que le titre en question.



