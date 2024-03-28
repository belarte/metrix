all: generate
	go build -o bin/

generate:
	templ generate

run: generate
	go run main.go serve

test: generate
	go test ./...  -fullpath

docker:
	docker compose up --build

migrate-prod:
	go run main.go db migrate -d bin/prod.sqlite

run-prod: migrate-prod docker
