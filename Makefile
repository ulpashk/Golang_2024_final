run/api:
	go run ./cmd/api


db/psql:
	psql ${DB_DSN}


db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext .sql -dir ./migrations ${name}


db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DB_DSN} up


db/migrations/down:
	@echo 'Running down migrations...'
	migrate -path ./migrations/000001 -database ${DB_DSN} down
