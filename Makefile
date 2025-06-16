default: go-run
.PHONY: go-run

MIGRATION_MODE_UP=up

PATH_DB_POSTGRES=postgres://evans:evans@localhost:8081/postgres?sslmode=disable
FILE_MIGRATIONS_POSTGRES =./migrations/postgres

PATH_DB_CLICKHOUSE=clickhouse://evans:evans@localhost:8083/logs
FILE_MIGRATIONS_CLICKHOUSE =./migrations/clickhouse

# Golang
go-run:
	go run cmd/main.go -config ./local.env

go-run-nats:
	go run cmd/events/events.go -config ./local.env

go-migrate-postgres-up:	
	go run cmd/migrator/migrator.go -mode $(MIGRATION_MODE_UP) -storage-path $(PATH_DB_POSTGRES) -migrations-path $(FILE_MIGRATIONS_POSTGRES)

go-migrate-clickhouse-up:	
	go run cmd/migrator/migrator.go -mode $(MIGRATION_MODE_UP) -storage-path $(PATH_DB_CLICKHOUSE) -migrations-path $(FILE_MIGRATIONS_CLICKHOUSE)

go-fmt:
	go fmt ./...

go-lint:
	golangci-lint run ./... -c .golangci.yml

go-memory-check:
	fieldalignment ./...

go-memory-fix:
	fieldalignment -fix ./...

go-mock:
	go generate ./...

go-cover:
	go test -cover ./...

go-swag:
	swag init -g cmd/main.go

# go-cover-html:
# 	go test -cover -coverprofile=coverage.out ./internal/users && go tool cover -html=coverage.out -o coverage.html


# cmd //c tree "%cd%" //F