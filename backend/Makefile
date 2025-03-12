.PHONY: build run docker docker-compose clean swagger docs

# Build the application
build:
	go build -o main .

# Run the application
run:
	go run main.go

# Build Docker image
docker:
	docker build -t lending-borrowing-api .

# Run the Docker Compose stack
docker-compose: env-setup
	docker-compose up --build

# Restart the Docker Compose stack
docker-restart: env-setup
	docker-compose down
	docker-compose up --build

# Stop the Docker Compose stack
docker-stop:
	docker-compose down

# Clean build artifacts
clean:
	rm -f main
	rm -rf tmp/

# Generate Swagger documentation
swagger:
	swag init -g main.go

# Run tests
test:
	go test -v ./...

# Tidy up dependencies
tidy:
	go mod tidy

# Install all project dependencies
deps:
	go mod download

# Create .env file from example
env-setup:
	test -f .env || cp .env.example .env
	echo "Don't forget to update values in .env file"

# Create initial database migration
migrate:
	go run main.go migrate

# Start development mode with hot reloading
dev:
	air

# Show all Docker containers
docker-ps:
	docker-compose ps

# View logs for the API service
logs:
	docker-compose logs -f app

# Open API docs in browser (MacOS)
docs:
	open http://localhost:8080/swagger/index.html
