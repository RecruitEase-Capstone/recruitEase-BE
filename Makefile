include .env

run:
	@go run cmd/main.go 

migrate-up:
	@migrate -path authentication_service/db/migrations -database "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up    

migrate-down:
	@migrate -path authentication_service/db/migrations -database "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

migrate-fix:
	@migrate -path authentication_service/db/migrations -database "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose force 20241106063649

compose-up:
	@docker-compose up -d --build

compose-down:
	@docker-compose down

auth-run:
	@go run authentication_service/cmd/main.go

.PHONY: run migrate-up migrate-down migrate-fix compose-up compose-down auth-run