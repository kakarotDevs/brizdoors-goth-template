# syntax = docker/dockerfile:1

# Build stage for Node.js (Tailwind CSS)
ARG NODE_VERSION=24.2.0
FROM node:${NODE_VERSION}-slim AS node-builder

WORKDIR /app

# Install packages needed to build node modules
RUN apt-get update -qq && \
    apt-get install --no-install-recommends -y build-essential node-gyp pkg-config python-is-python3

# Install node modules
COPY package-lock.json package.json ./
RUN npm ci

# Copy application code and build CSS
COPY . .
RUN npm run build:css

# Build stage for Go
FROM golang:1.21-alpine AS go-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from go-builder stage
COPY --from=go-builder /app/app .

# Copy the built CSS from node-builder stage
COPY --from=node-builder /app/public/styles.css ./public/styles.css

# Copy other static files
COPY --from=node-builder /app/public/branding ./public/branding

# Expose port
EXPOSE 3000

# Run the binary
CMD ["./app"]
