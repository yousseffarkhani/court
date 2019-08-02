TODOs :

- Add login and register mechanism
- Commentaires
- Copier le README du projet NODE
- Remplacer "github.com/yousseffarkhani/court/courtdb"
  "github.com/yousseffarkhani/court/views" par views / courtdb / ...

# Introduction

# Commands

Create .env file with 3 variables (APP_ENV, POSTGRES_USER, POSTGRES_PASSWORD)

- Launch project in dev mode :

1. Delete APP_ENV from .env file.
2. `docker-compose up --build`

- Launch project in production mode :

1. Add APP_ENV=production to .env file.
2. `docker-compose up --build`
