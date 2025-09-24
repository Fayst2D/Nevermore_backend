build:
	go build -o app main.go

run:
	go run main.go

set_dbstring:
	export GOOSE_DBSTRING="user=postgres password=1 dbname=userservice sslmode=disable" && echo $$GOOSE_DBSTRING

set_driver:
	export GOOSE_DRIVER=postgres

clean:
	rm -f app

test_cover:
	go test -cover ./...


