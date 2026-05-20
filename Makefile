include .env

# Migration commands
migrate_up:
	goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)" up

migrate_down:
	goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)" down

migrate_status:
	goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)" status

.PHONY: migrate_up migrate_down migrate_status