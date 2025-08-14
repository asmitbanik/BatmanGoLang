.PHONY: build run test docker-up docker-build lint

build:
	go build -o shazam ./main.go

run: build
	./shazam

test:
	go test ./... -v

docker-up:
	docker-compose up --build

docker-build:
	docker build -t shazam-for-code:latest .
