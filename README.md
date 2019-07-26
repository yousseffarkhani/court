X Supprimer et réinitialiser le git

- Faire en sorte qu'au début de l'application la BDD soit initialisée.
- Copier le README du projet NODE
- Ajouter le fichier .env pour masquer les passwords

# Introduction

# Commands

Create .env file with 3 variables (APP_ENV, POSTGRES_USER, POSTGRES_PASSWORD)

- Launch project in dev mode :

1. Delete APP_ENV from .env file.
2. `docker-compose up --build`

- Launch project in production mode :

1. Add APP_ENV=production to .env file.
2. `docker-compose up --build`
