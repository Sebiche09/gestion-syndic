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
- Installer devpod sur la machine : https://devpod.sh/
- cliquer sur le lien : [![Ouvrir dans devPod!](https://devpod.sh/assets/open-in-devpod.svg)](https://devpod.sh/open#https://github.com/Sebiche09/gestion-syndic)
- crée le workspace en donnant un nom à l'espace de travail
exemple :
![image](https://github.com/user-attachments/assets/d76bbb6e-cf35-4b67-830e-fb7176af89b1)
- une fois vscode de lancer, si une erreur apparait, simplement faire more actions... fermer vscode et le relancer.
- une fois dans le terminal , tapper la commande docker-compose up
Le projet est pret
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
