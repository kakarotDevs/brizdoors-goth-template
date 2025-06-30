.DEFAULT_GOAL := build

# Build the application
build:
	@echo "Building application..."
	if not exist bin mkdir bin
	@templ generate
	@go build -o bin/app .

# Run the built application
run: build
	@echo "Starting application..."
	@./bin/app

# Development mode with hot reload
dev:
	@echo "Starting development environment..."
	npx concurrently \
		"npx tailwindcss -i views/css/app.css -o public/styles.css --watch" \
		"templ generate --watch --proxy=http://localhost:3000" \
		"air"

# Build CSS only
build-css:
	@echo "Building CSS..."
	npx tailwindcss -i views/css/app.css -o public/styles.css

# Generate templates only
generate:
	@echo "Generating templates..."
	templ generate

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin
	@rm -f public/styles.css
	@find . -name "*_templ.go" -delete 2>/dev/null || true

# Install dependencies
install:
	@echo "Installing dependencies..."
	@npm install
	@go mod tidy

# Help target
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  run        - Build and run the application"
	@echo "  dev        - Start development environment with hot reload"
	@echo "  build-css  - Build CSS only"
	@echo "  generate   - Generate templates only"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install dependencies"
	@echo "  help       - Show this help message"

.PHONY: build run dev build-css generate clean install help
