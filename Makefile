.PHONY: test test-python test-go test-ts lint up down build clean

# Run all tests
test: test-python test-go test-ts

# Python tests
test-python:
	cd services/analyzer && pip install -q -r requirements.txt && pytest -v

# Go tests
test-go:
	cd services/gateway && go test ./... -v

# TypeScript tests
test-ts:
	cd services/dashboard && npm install --silent && npm test

# Run all linters
lint: lint-python lint-go lint-ts

lint-python:
	cd services/analyzer && flake8 app/ tests/ --max-line-length=120

lint-go:
	cd services/gateway && go vet ./...

lint-ts:
	cd services/dashboard && npx eslint 'src/**/*.ts'

# Docker Compose commands
up:
	docker compose up -d --build

down:
	docker compose down

build:
	docker compose build

# Health check all services
health:
	@echo "Analyzer:" && curl -sf http://localhost:8001/health | python3 -m json.tool || echo "  unavailable"
	@echo "Gateway:" && curl -sf http://localhost:8002/health | python3 -m json.tool || echo "  unavailable"
	@echo "Dashboard:" && curl -sf http://localhost:8003/health | python3 -m json.tool || echo "  unavailable"

clean:
	docker compose down -v --rmi local
	find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
	find . -type d -name node_modules -exec rm -rf {} + 2>/dev/null || true
	find . -type d -name dist -exec rm -rf {} + 2>/dev/null || true
