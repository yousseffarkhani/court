# TODOs :
- Testing
- Amélioration du code
- Finir les TODOs
- Déployer
- Voir https://data.iledefrance.fr/explore/dataset/20170419_res_fichesequipementsactivites/information/?disjunctive.actlib

# Commandes

Create .env file with 2 variables (POSTGRES_USER, POSTGRES_PASSWORD)

- Launch project in dev mode :

1. `docker-compose up --build`

- Launch project in production mode :

1. Add APP_ENV=production to .env file.
2. `docker-compose up --build`

# Description de l'application

## Introduction

Court est une application permettant de trouver le terrain de basket le plus proche dans Paris (pour le moment).
Pour cela, les données ont été récupérées depuis plusieurs sources (Open data, web scraping).
Il est possible de créer un compte permettant de commenter les terrains et d'en soumettre de nouveaux.

## Features

### Web Scraper (NodeJS, librairie : puppeteer/cheerio)

- Data scraper (DONE)

### Application :

- Liste des terrains parisiens (DONE)
- Notation des terrains
- Commentaires (DONE)
- Responsive (DONE)
- Niveau de jeu des terrains
- Localisation des terrains (Gmap) (DONE)
- Description (DONE)
- Mise en place de filtres
- PWA
- Soumettre de nouveaux terrains et créer une page admin pour les accepter
- Photos
- Réconciliation de données
- Recherche par arrondissement (DONE)
- Agenda des terrains et création de communautés

## Intention

- Mettre en oeuvre les connaissances apprises en Golang, NodeJS, docker, HTML et CSS.
- Réaliser le déploiement d'un site
- Réaliser un site utile à terme :)

## Remarques

- Le projet back-end mélange plusieurs types d'architectures :

  - API appelées depuis les pages HTML
  - Peuplement des données directement dans les pages HTML. Cela a été fait de manière intentionnelle afin de tester plusieurs concepts. Il se peut qu'il en résulte une application non optimisée.

-Pas assez de temps pour mettre en place des tests.

- Réduire au minimum le javascript utilisé pour économiser les ressources du client, simplifier le code, améliorer la performance du crawler Google.

## Réalisation

1. Web scraping des données
   - Récupération des données des terrains
     - Nom
     - Adresse
     - Lat/long
     - Dimensions
     - Revetement
     - Découvert ?
     - éclairage ?
2. Mise en place du back-end
   - Enregistrement des données dans une BDD

## Sources de données

- https://www.gralon.net/mairies-france/paris/equipements-sportifs-terrain-de-basket-ball-75056.htm
- http://www.cartes-2-france.com/activites/750560006/ritz-health-club.php donne accès aux liens https://www.webvilles.net/sports/75056-paris.php
- https://www.data.gouv.fr/fr/datasets/recensement-des-equipements-sportifs-espaces-et-sites-de-pratiques/
