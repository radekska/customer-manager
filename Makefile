migrate:
	go run cmd/migrate/migrate.go

start:
	go run cmd/server/server.go

tests:
	go test -v ./...

api-docs:
	swag init