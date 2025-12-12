# Multi-stage build for optimal image size

# Stage 1: Builder
FROM golang:1.24.3-alpine AS builder

ARG VERSION=dev
ARG BUILD_TIME=unknown

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -a -installsuffix cgo \
    -o copilot-os ./cmd/server

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/copilot-os .

# Create non-root user for security
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

USER app

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/app/copilot-os", "health"] || exit 1

EXPOSE 8080

ENTRYPOINT ["/app/copilot-os"]
CMD ["serve"]
