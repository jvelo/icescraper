build:
	go build -o icescraper main.go

migrate:
	. .env && npx prisma migrate dev

generate:
	go generate ./...

test:
	go test ./...

run:
	go run main.go

lint:
	golangci-lint run

docker:
	docker build -t jvelo/icescraper .

all: build
