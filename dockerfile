## Build
FROM golang:1.20.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o docker-app ./cmd/main.go

CMD ["./docker-app"]

## Run
FROM alpine:latest

WORKDIR /app

COPY .env .

COPY --from=builder /app/docker-app .

CMD ["./docker-app"]

## docker build -t docker-app .