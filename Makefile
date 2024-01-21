all:
	go build -o bin/

run:
	go run main.go serve

test:
	go test ./...  -fullpath

docker:
	docker-compose up --build
