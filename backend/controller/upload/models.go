package upload

// OwnerInfo stocke les informations d'un propriétaire
type OwnerInfo struct {
	LastName  string      `json:"last_name"`
	FirstName string      `json:"first_name"`
	Address   AddressInfo `json:"address"`
	Title     string      `json:"title"`
}

// AddressInfo stocke les informations d'adresse
type AddressInfo struct {
	Street     string `json:"street"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
	Country    string `json:"country"`
}

// PostalInfo représente les informations récupérées depuis l'API, comme le code postal.
type PostalInfo struct {
	PostalCode string `json:"code_postal"`
}
