.PHONY: migrate
migrate:
	goose -dir=internal/migrations postgres "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable" up

downgrade:
	goose -dir=internal/migrations postgres "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable" down

.PHONY: db
db:
	docker run -d --name workout-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres -p 5432:5432 postgres:16

.PHONY: redis
redis:
	docker run -d --name workout-redis -p 6379:6379 redis:latest

.PHONY: run
run:
	go run main.go
