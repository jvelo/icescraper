build:
	go build -o icecast-monitor main.go

migrate:
	go get github.com/prisma/prisma-client-go@v0.4.0
	go run github.com/prisma/prisma-client-go db push --preview-feature

generate:
	go get github.com/prisma/prisma-client-go@v0.4.0
	go run github.com/prisma/prisma-client-go generate

run:
	go run main.go

lint:
	golangci-lint run

all: build
