FROM golang:1.16 as build

WORKDIR /app

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# build the binary with all dependencies
RUN go build -o /icescraper .

CMD ["/icescraper"]