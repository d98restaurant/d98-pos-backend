# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pos-backend cmd/api/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/pos-backend .

# Copy .env file if needed (optional)
COPY --from=builder /app/.env .env

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./pos-backend"]
