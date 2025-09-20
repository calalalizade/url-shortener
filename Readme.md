# Start

docker compose down -v
docker compose build
docker compose up -d db redis
docker compose run --rm migrate
docker compose up -d api

# Rollback

docker compose run --rm migrate down

# Check health

docker compose logs -f db
docker compose exec redis redis-cli ping
