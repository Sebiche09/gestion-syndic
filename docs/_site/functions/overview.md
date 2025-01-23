# Vue d'ensemble des fonctions

Ce document fournit une vue d'ensemble des fonctions implémentées dans le backend et le frontend du projet. Il détaille leur rôle principal, leur structure, et leurs interactions.

---

## **Fonctions Backend**
Le backend gère la logique métier, les interactions avec la base de données, et expose des endpoints via une API REST.

### 1. [condominium.go]({{ site.baseurl }}/functions/backend/condominium)
- **Description :** Contient les fonctions liées à la gestion des condominiums.
- **Fonctions principales :**
  - `GetAllCondominiums` : Récupère tous les condominiums avec leurs adresses associées.
  - `CreateCondominium` : Crée un nouveau condominium avec ses occupants et ses unités.
  - `CheckUniqueness` : Vérifie si un nom ou un préfixe est unique dans la base de données.
  - `ChecjIfExists` : Vérifie si une valeur existe déjà dans une table donnée.
  - `getOccupantTypeLabel` : Retourne le label correspondant à une abréviation trouvée dans un titre.


### 2. [unit.go]({{ site.baseurl }}/functions/backend/unit)
- **Description :** Gère les unités associées aux condominiums.
- **Fonctions principales :**
  - `GetUnits` : Retourne les informations détaillées sur toutes les unités.
  - `CreateUnit` : Ajoute une nouvelle unité au condominium.

### 3. [occupant.go]({{ site.baseurl }}/functions/backend/occupant)
- **Description :** Traite les occupants associés aux condominiums.
- **Fonctions principales :**
  - `GetOccupants` : Liste les occupants associés à un condominium.
  - `AddOccupant` : Ajoute un nouvel occupant avec ses informations personnelles et professionnelles.
