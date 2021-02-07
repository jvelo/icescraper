FROM golang:1.15 as build

WORKDIR /app

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/prisma/prisma-client-go prefetch

COPY . ./

# generate the Prisma Client Go client
RUN go generate ./...

# build the binary with all dependencies
RUN go build -o /icecast-monitor .

CMD ["/icecast-monitor"]