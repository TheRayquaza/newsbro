help:
    just --list

# Init db
init-db:
    #!/bin/sh
    docker compose up postgres
    cd apps && docker compose exec postgres bash -c "
        psql -U username -d postgres -c 'CREATE DATABASE account;';
        psql -U username -d postgres -c 'CREATE DATABASE article;';
        psql -U username -d postgres -c 'CREATE DATABASE analytics;';
    "

# Runs all servers
dev:
    cd apps/ && docker compose up -d
    cd apps/port-front/ && npm i && npm run dev
