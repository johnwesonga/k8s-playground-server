# ---- Build Stage ----
FROM golang:1.25.3-alpine AS builder

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

WORKDIR /app

# Install git (needed for some dependencies)
RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o server

# ---- Runtime Stage ----
FROM alpine:3.23

WORKDIR /app

# Add non-root user
RUN adduser -D appuser

# Copy binary from builder
COPY --from=builder /app/server .

USER appuser

EXPOSE 8080

CMD ["./server"]