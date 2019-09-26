La partie back-end est en pleine refonte (cf. https://github.com/yousseffarkhani/Playground-Back-end).
N'étant pas satisfait des capacités du templating Golang, je vais ajouter un front-end que je réaliserai à l'aide d'un des 3 frameworks (Angular/React/Vue).
# TODOs :
- Vérification de la casse des noms de terrains au moment de les ajouter
- Ajouter du testing
- Refactorer le code (DRY, utiliser plus d'interfaces, changer la structure du projet)
- Finir les TODOs dans le code
- Testing
- Amélioration du code
- Homogénéiser la langue utilisée (français)
- ~~Ajouter une fonction permettant de lancer le serveur en https en mode production et sur le port 8080 de localhost en mode dev.~~
- Rediriger le traffic HTTP vers HTTPS
- Automatiser le renouvellement du certificat TLS.


# Améliorations
- Améliorer le front-end
- Ajouter d'autres jeux de données : https://data.iledefrance.fr/explore/dataset/20170419_res_fichesequipementsactivites/information/?disjunctive.actlib

# Commandes

Créer un fichier .env dans la racine du projet avec 2 variables (POSTGRES_USER=XXX, POSTGRES_PASSWORD=XXX)

- Lancer le projet en mode dev (hot reload) :

1. `docker-compose up --build`

- Lancer le projet en mode production sur AWS :

1. Ajouter APP_ENV=production au fichier .env.
2. Ajouter
```
CERTFILE=XXX/cert.pem
PRIVKEY=XXX/privkey.pem

```
3. Changer dans le Dockerfile et docker-compose les ports sur 443 pour du HTTPS.
4. `docker-compose up --build`

# Description de l'application

## Introduction

Court est une application permettant de trouver le terrain de basket le plus proche dans Paris.
Pour cela, les données ont été récupérées depuis plusieurs sources (Open data, web scraping).
Il est possible de créer un compte permettant de commenter les terrains et d'en soumettre de nouveaux.

## Fonctionnalités
Le webscraper produit un fichier JSON qui sert ensuite à peupler la base de données de l'application.
### Web Scraper (NodeJS, librairie : puppeteer/cheerio)

- ~~Data scraper~~

### Application :

- ~~Liste des terrains parisiens~~
- Notation des terrains
- ~~Commentaires~~
- ~~Responsive~~
- Niveau de jeu des terrains
- ~~Localisation des terrains (Gmap)~~
- ~~Description des terrains~~
- Page profil de l'utilisateur
- ~~Utilisation de JWT pour garder la session active~~
- ~~Utilisation de PostgreSQL pour enregistrer les utilisateurs, terrains et commentaires~~
- Mise en place de filtres
- PWA
- ~~Soumettre de nouveaux terrains~~ et créer une page admin pour les accepter
- Photos
- Réconciliation de données
- ~~Recherche par arrondissement~~
- Agenda des terrains et création de communautés
- ~~Déployer l'application en https~~

## Intention

- Mettre en oeuvre les connaissances apprises en Golang, NodeJS, docker, HTML et CSS.
- Réaliser le déploiement d'un site en HTTPS avec un nom de domaine acheté
- Réaliser un site utile à terme :-)

## Remarques

- Le projet back-end mélange plusieurs types d'architectures :

  - API appelées depuis les pages HTML
  - Peuplement des données directement dans les pages HTML. Cela a été fait de manière intentionnelle afin de tester plusieurs concepts. Il se peut qu'il en résulte une application non optimisée.

- Je n'ai pas démarré le projet en TDD. De ce fait, je n'ai pas eu assez de temps pour mettre en place des tests. Cependant, je suis en train de refaire la partie back-end en TDD.

- Le but recherché en utilisant le templating Golang était de réduire au minimum le javascript utilisé pour économiser les ressources du client, simplifier le code, améliorer la performance du crawler Google.

## Sources de données

- https://www.gralon.net/mairies-france/paris/equipements-sportifs-terrain-de-basket-ball-75056.htm
- http://www.cartes-2-france.com/activites/750560006/ritz-health-club.php donne accès aux liens https://www.webvilles.net/sports/75056-paris.php
- https://www.data.gouv.fr/fr/datasets/recensement-des-equipements-sportifs-espaces-et-sites-de-pratiques/
