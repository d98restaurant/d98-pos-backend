.PHONY: run build clean test install

run:
	go run cmd/api/main.go

build:
	go build -o bin/pos-backend cmd/api/main.go

clean:
	rm -rf bin/

test:
	go test -v ./...

install:
	go mod download
	go mod tidy

docker-build:
	docker build -t pos-backend .

docker-run:
	docker run -p 8080:8080 --env-file .env pos-backend