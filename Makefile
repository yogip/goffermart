
mock:
	mockgen -destination=internal/mocks/mock_system_service.go -package=mocks metrics/internal/core/service Pinger	

test:
	go test ./...

server:
	go run cmd/gophermart/main.go -l debug -d "host=localhost port=25432 user=username password=password dbname=gophermart sslmode=disable"

migrate-new:
	goose create $@ sql

migrate-up:
	goose -dir migrations postgres "user=username dbname=gophermart password=password sslmode=disable host=127.0.0.1 port=25432" up

migrate-down:
	goose -dir migrations postgres "user=username dbname=gophermart password=password sslmode=disable host=127.0.0.1 port=25432" down