build:
	go build -o app main.go

run:
	go run main.go

clean:
	rm -f app

test_cover:
	go test -cover ./...

migration_create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migration_create name=<migration_name>"; \
		exit 1; \
	fi
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=1 dbname=nevermore sslmode=disable" goose create $(name) sql

goose_up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=1 dbname=nevermore sslmode=disable host=localhost" goose up
goose_down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=1 dbname=nevermore sslmode=disable host=localhost" goose down