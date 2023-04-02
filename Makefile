migrate:
	go run cmd/migrate/migrate.go

start:
	go run cmd/server/server.go

api-docs:
	swag init