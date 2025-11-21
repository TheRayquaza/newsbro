# Shows this help
help:
    just --list

# Inits database table (run this the first time you launch the project)
init-db:
    #!/bin/sh
    cd apps
    docker compose up -d postgres
    docker compose exec postgres bash -c "
        psql -U username -d postgres -c 'CREATE DATABASE account;';
        psql -U username -d postgres -c 'CREATE DATABASE article;';
        psql -U username -d postgres -c 'CREATE DATABASE analytics;';
    "

# Runs all servers & frontend using docker compose and npm
dev:
    cd apps/ && docker compose up -d --build
    cd apps/port-front/ && npm i && npm run dev

# Shuts down docker compose
down:
    cd apps/ && docker compose down
