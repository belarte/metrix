all:
	go build -o bin/

generate:
	templ generate

run: generate
	go run main.go serve

test: generate
	go test ./...  -fullpath

docker:
	docker-compose up --build
