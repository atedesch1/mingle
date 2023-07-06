build:
	go build -o bin/mingle

run: build
	DATABASE_URL="postgresql://root:pass@0.0.0.0:5432/mingle?sslmode=disable" ./bin/mingle

test:
	go test ./...

postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pass -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root mingle

dropdb:
	docker exec -it postgres15 dropdb mingle

migrateup:
	migrate -path db/migrations -database "postgresql://root:pass@0.0.0.0:5432/mingle?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:pass@0.0.0.0:5432/mingle?sslmode=disable" -verbose down

.PHONY: 
	createdb
	dropdb
