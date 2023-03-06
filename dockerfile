# Build stage
FROM golang:1.20.1-alpine3.17 AS build

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download the dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# Final stage
FROM alpine:3.17.2

# Set the working directory
WORKDIR /app

# Copy the config file
COPY .env .

# Copy the binary from the build stage
COPY --from=build /app/app .

# Run the binary
CMD ["./app"]
