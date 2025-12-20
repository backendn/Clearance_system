postgres:
	docker run --name postgres19 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine
createdb:
	docker exec -it postgres19 createdb --username=root --owner=root university_clearance
dropdb:
	docker exec -it postgres19 dropdb university_clearance
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/university_clearance?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/university_clearance?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/backendn/simplebank/db/sqlc Store

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock