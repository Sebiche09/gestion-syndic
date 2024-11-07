# gestion-syndic
Automatisation de la gestion des immeubles pour les syndics.

### Politique de commit
a.b.c - détail du commit
La version est décompoée en trois parties :
- a => version majeur
- b => version mineur
- c => patch de version

Par exemple :
- la version avec le scan de facture fonctionnel serait nomée : 0.1.0 - scan de facture fonctionnel
- la version avec les templates serait nomée : 0.2.0 - templates de facture fonctionnel

## Fonctionnalités à dev
- scan de factures
- utilisation de templates de facture
- création des templates de facture
- scan des extrets de compte
- gestion du compte banquaire
- possibilité d'enregistrer plusieurs batiments différents
- gestion des comptes différents

## Infrastructure
Dans un premier temps, il faudra utiliser un serveur chez Hostinger
- container DB + agent zabbix
- container web + agent zabbix
- container zabbix
- Serveur isolé pour backups

## Lancement du programme : 
Dans le dossier principale gestion-syndic, faire la commande : 
- docker-compose up -d
- docker exec -it ollama ollama run llama3
- /exit

## Code d'erreur : 

Le code d'erreur est divisé en 3 parties
le code commence par la lettre E
Le premier chiffre désigne si la catégorie d'erreur (go - 1, angular - 2, base de donnée - 3). 
Le deuxième chiffre désigne la nature de l'erreur : 
| code erreur | nature |
| ----------- | ------ |
| 1           | erreur de format|
| 2           | donnée deja existante|
| 3           | Problème d'accès à la db|
|4            | Problème avec une fonction de golang|

devpod up https://github.com/Sebiche09/gestion-syndic
