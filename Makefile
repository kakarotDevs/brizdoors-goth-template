.DEFAULT_GOAL := build

run: build
	@./bin/app

build:
	@mkdir -p bin
	@go build -o bin/app .

dev:
	@echo "Starting dev environment..."
	npx concurrently \
		"npx tailwindcss -i views/css/app.css -o public/styles.css --watch" \
		"templ generate --watch --proxy=http://localhost:3000" \
		"air"

clean:
	@rm -rf bin

run-dev: dev
