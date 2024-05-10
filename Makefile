include .env

postgres:
	docker run --name postgres_serenify -p ${DATABASE_PORT}:${DATABASE_PORT} -e POSTGRES_USER=${DATABASE_USER} -e POSTGRES_PASSWORD=${DATABASE_PASSWORD} -d postgres:16-alpine3.19

createdb:
	docker exec -it postgres_serenify createdb --username=${DATABASE_USER} --owner=${DATABASE_USER} ${DATABASE_NAME}

dropdb:
	docker exec -it postgres_serenify dropdb ${DATABASE_NAME}

migrateup:
	migrate -path migration -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose up

migratedown:
	migrate -path migration -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose down

migrateforce:
	migrate -path migration -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose force 20240508120431

server:
	go run cmd/main.go

.PHONY : migrateup migratedown server migrateforce