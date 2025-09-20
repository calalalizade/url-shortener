docker compose down -v

docker compose build

docker compose up -d db

docker compose run --rm migrate

docker compose up -d api

- Rollback
  docker compose run --rm migrate down
