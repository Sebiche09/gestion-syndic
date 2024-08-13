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

##Lancement du programme : 
Dans le dossier principale gestion-syndic, faire la commande : docker-compose up -d
une fois fini, 
faire la commande : docker exec -it ollama ollama run llama3
une fois fini,
faire la commande /exit

