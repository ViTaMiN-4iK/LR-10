# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go-service/go.mod go-service/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code (сохраняем структуру папок)
COPY go-service/ .

# Build the application
RUN go build -o main .

# Stage 2: Run
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080 50051

CMD ["./main"]