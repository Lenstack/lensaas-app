## Build
FROM golang:1.20.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/*.go ./

RUN go build -o /docker-app

## Deploy
FROM alpine:latest

WORKDIR /

COPY --from=build ./docker-app ./docker-app

EXPOSE 8080

ENTRYPOINT ["/docker-app"]

## docker build -t docker-app .