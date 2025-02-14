postgres:
	 docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	 docker exec -it postgres12 createdb --username=root --owner=root wallet_bank

dropdb:
	 docker exec -it postgres12 dropdb wallet_bank

migrateup: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/wallet_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/wallet_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock_db:
	mockgen -package mockdb -destination db/mock/store.go github.com/9Neechan/JavaCode-test-task/db/sqlc Store

mock_rabbitmq:
	mockgen -package mockrabbitmq -destination rabbitmq/mock/rabbitmq.go github.com/9Neechan/JavaCode-test-task/rabbitmq AMQPChannel

build: postgres createdb migrateup server

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock build
