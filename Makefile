include .env

migrate_up:
	migrate -database ${PSQL_MIGRATION_URL} -path pkg/db/migration up

migrate_down:
	migrate -database ${PSQL_MIGRATION_URL} -path pkg/db/migration  down 1