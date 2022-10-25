postgres:
	docker run --name postgres  -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres

createdb:
	docker exec -it postgres createdb --username=postgres sample_bank

dropdb:
	docker exec -it postgres dropdb --username=postgres sample_bank

migrateup:
	migrate -path db/migrations/up -database 'postgresql://postgres:postgres@localhost:5432/sample_bank?sslmode=disable' up

migratedown:
	migrate -path db/migrations/down -database 'postgresql://postgres:postgres@localhost:5432/sample_bank?sslmode=disable' down

sqlc:
	sqlc generate
test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test