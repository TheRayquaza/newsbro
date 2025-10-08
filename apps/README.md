# Apps

You can start apps using the docker compose provided

```bash
docker compose up -d
```

Then start the apps and the frontend:

```bash
cd port-front && npm run dev # one term (port 5173)
cd repo-account && go build -o bin/api src/cmd/main.go && ./bin/api # 2nd (port 8080)
cd repo-article && go build -o bin/api src/cmd/main.go && PORT=8081 ./bin/api # 3rd (port 8081)
## others ?
```

Inspect db:

```bash
docker exec -it apps-postgres-1 bash
psql -U username -d postgres
\c your_db
```