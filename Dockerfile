# Multi-stage build for WhatsApp Go MCP
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy source code
COPY src/ ./src/

# Download dependencies
WORKDIR /app/src
RUN go mod download

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o whatsapp .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    nodejs \
    npm

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/src/whatsapp .

# Copy static files and views
COPY src/statics ./statics
COPY src/views ./views

# Create necessary directories including storages for persistent data
RUN mkdir -p .wwebjs_auth statics/media statics/qrcode statics/senditems storages

# Create volume for persistent WhatsApp session and chat data
VOLUME ["/app/storages"]

# Set environment variables
ENV PORT=8081 \
    PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

# Expose port
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run MCP server
CMD ["./whatsapp", "mcp"]