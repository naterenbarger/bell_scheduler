.PHONY: install test test-coverage build serve clean lint

# Install dependencies
install:
	npm install

# Run unit tests
test:
	npm run test:unit

# Run tests with coverage
test-coverage:
	npm run test:unit -- --coverage

# Build for production
build:
	npm run build

# Serve development server
serve:
	npm run serve

# Clean build artifacts
clean:
	rm -rf dist/
	rm -rf coverage/
	rm -rf node_modules/

# Run linter
lint:
	npm run lint

# Run all checks (lint, test, build)
check: lint test build

# Run tests in watch mode
test-watch:
	npm run test:unit -- --watch

# Run tests with UI
test-ui:
	npm run test:unit -- --ui 