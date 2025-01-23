# Table des matières
- [Aperçu des fonctions]({{ site.baseurl }}/functions/overview)
- [Documentation API]({{ site.baseurl }}/api)

# Introduction à Gestion-Syndic

**Gestion-Syndic** est une application web complète permettant de gérer efficacement les syndicats de copropriété. Ce projet repose sur une architecture modulaire et scalable, combinant un frontend moderne avec Angular et un backend robuste développé en Go.

---

## Objectifs du projet

Le projet vise à :
1. **Simplifier la gestion des syndicats** :
   - Gestion des occupants et des copropriétaires.
   - Gestion des unités (logements, bureaux, etc.) et des informations associées.
2. **Automatiser les tâches** :
   - Traitement des documents grâce à l'intégration d'une solution OCR.
   - Chargement et gestion de fichiers via FTP sécurisé.
3. **Fournir un environnement collaboratif** :
   - Interface web intuitive.
   - API backend permettant des intégrations tierces.

---

## Fonctionnalités principales

1. **Gestion des données** :
   - Création, modification et suppression des copropriétaires, unités, et informations associées.
   - Organisation des informations liées aux immeubles.

2. **OCR et traitement des fichiers** :
   - Extraction automatique des données à partir de documents téléversés.
   - Gestion des fichiers via un serveur FTP intégré.

3. **Scalabilité et sécurité** :
   - Architecture conteneurisée avec Docker pour faciliter le déploiement.
   - Sécurisation des accès aux données sensibles.

4. **Interface utilisateur** :
   - Application Angular intuitive avec navigation optimisée.
   - Visualisation des données en temps réel.

---

## Architecture du projet

### Structure générale
Le projet suit une architecture en trois tiers :
1. **Frontend** :
   - Développé avec Angular.
   - Propose une interface utilisateur moderne et responsive.
2. **Backend** :
   - Écrit en Go pour garantir des performances élevées.
   - Gère la logique métier et les connexions à la base de données PostgreSQL.
3. **Services complémentaires** :
   - Serveur FTP pour le stockage des fichiers.
   - OCR intégré pour automatiser le traitement des documents.

### Organisation des conteneurs
Le projet utilise Docker et `docker-compose` pour orchestrer les services :
- **PostgreSQL** : Base de données relationnelle.
- **Angular** : Frontend hébergé sur un serveur Node.js.
- **Go API** : Backend exposant des API REST.
- **OCR** : Service de traitement des documents.
- **FTP** : Serveur de gestion des fichiers.

---

## Technologies utilisées

- **Frontend** : Angular (TypeScript, SCSS).
- **Backend** : Go (Golang).
- **Base de données** : PostgreSQL.
- **Infrastructure** : Docker, Docker Compose.
- **Autres** :
  - Compodoc pour la documentation Angular.
  - Air pour le rechargement à chaud du backend.

---

## À qui s'adresse ce projet ?

Ce projet est conçu pour :
- Les gestionnaires de syndic de copropriété cherchant une solution digitale.
- Les développeurs souhaitant contribuer à une application modulaire.
- Toute organisation cherchant à automatiser la gestion des documents et des données liées aux copropriétés.

---

## Ressources supplémentaires

- [GitHub officiel](https://github.com/Sebiche09/gestion-syndic)