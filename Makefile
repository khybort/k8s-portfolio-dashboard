.PHONY: help build up down logs test clean

help:
	@echo "Available commands:"
	@echo "  make up       - Start all services with docker-compose"
	@echo "  make down     - Stop all services"
	@echo "  make logs     - Show logs from all services"
	@echo "  make build    - Build all Docker images"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean up containers and volumes"

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

build:
	docker-compose build

test:
	cd backend && go test ./...
	cd auth-service && go test ./...

clean:
	docker-compose down -v
	docker system prune -f

