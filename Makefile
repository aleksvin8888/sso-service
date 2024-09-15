# Makefile

run-local:
	go run cmd/sso/main.go --config=./config/local.yaml

run-migration:
	go run ./cmd/migrator/main.go --storage-path=./storage/sso.db --migrations-path=./migrations

run-test-migrate:
	go run ./cmd/migrator/main.go --storage-path=./storage/sso.db --migrations-path=./tests/migrations --migrations-table=migrations_test