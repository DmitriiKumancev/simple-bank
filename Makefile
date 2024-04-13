# Эти команды полезны чтобы легко настроить проект на локальном компьютере для дальнейшего развития
postgres:
	docker run --name simplebank-postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d -p 5438:5432 postgres:16-alpine

createdb:
	docker exec -it simplebank-postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simplebank-postgres dropdb simple_bank

.PHONY: postgres createdb dropdb