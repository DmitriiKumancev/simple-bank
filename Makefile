# Эти команды полезны чтобы легко настроить проект на локальном компьютере для дальнейшего развития
postgres:
	docker run --name simplebank-postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d -p 5438:5432 postgres:16-alpine

createdb:
	docker exec -it simplebank-postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simplebank-postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5438/simple_bank?sslmode=disable" -verbose up
	
migratedown:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5438/simple_bank?sslmode=disable" -verbose down
	
sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc