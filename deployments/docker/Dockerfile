# Multi-stage build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy source code
COPY . .

# Build the daemon
RUN make build

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -g 1000 daemon && adduser -D -u 1000 -G daemon daemon

# Copy binary from builder
COPY --from=builder /app/build/mydaemon /usr/local/bin/
COPY --from=builder /app/configs/mydaemon.yaml /etc/mydaemon/

# Create log directory
RUN mkdir -p /var/log/mydaemon && chown daemon:daemon /var/log/mydaemon

# Switch to non-root user
USER daemon

# Run daemon
CMD ["mydaemon", "-config", "/etc/mydaemon/mydaemon.yaml"]
