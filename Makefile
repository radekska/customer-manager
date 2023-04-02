migrate:
	go run cmd/migrate/migrate.go

start:
	go run cmd/server/server.go

docs:
	swag init