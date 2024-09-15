# Makefile

run-local:
	go run cmd/sso/main.go --config=./config/local.yaml

run-migration:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations
