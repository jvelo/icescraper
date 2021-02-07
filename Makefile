build:
	go build -o icecast-monitor main.go

migrate:
	go run github.com/prisma/prisma-client-go db push --preview-feature

generate:
	go generate ./...

test:
	DATABASE_URL=file:./test.db go run github.com/prisma/prisma-client-go db push --preview-feature
	go test ./...

run:
	go run main.go

lint:
	golangci-lint run

docker:
	docker build -t jvelo/icecast-monitor .

all: build
