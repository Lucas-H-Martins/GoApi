run-dev:
	GO_ENV=dev go run ./cmd/goapi/main.go

run-migrate-dev-up:
	GO_ENV=dev go run ./cmd/migrate/main.go -up

run-migrate-dev-down:
	GO_ENV=dev go run ./cmd/migrate/main.go -down
