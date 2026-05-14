# ============================================================
# Stage 1: Builder
# ============================================================
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache dependencies first (better layer caching)
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-w -s -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /app/bin/fish-game-api \
    ./cmd/api/main.go

# ============================================================
# Stage 2: Final (minimal runtime image)
# ============================================================
FROM alpine:3.19

# Security: run as non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Add CA certs and timezone data for HTTPS + time operations
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy compiled binary from builder
COPY --from=builder /app/bin/fish-game-api .

# Copy .env if it exists (optional, prefer env vars in docker-compose)
COPY --chown=appuser:appgroup .env.example .env.example

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/api/v1/health || exit 1

ENTRYPOINT ["./fish-game-api"]
