FROM golang:1.22-alpine3.18 AS builder

WORKDIR /build

COPY . .

RUN go build -o main ./cmd/api/main.go

CMD ["/build/main"]


