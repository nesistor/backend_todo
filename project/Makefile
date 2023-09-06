TASK_BINARY=taskApp

## up: starts all containers in the background without forcing build
up: 
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-cmpose (if running), builds all projects and start docker compose
up_build: build_broker-service
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down: 
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_task-service:
	@echo "Building task-service binary..."
	cd ../task-service && env GOOS=linux CGO=ENABLED=0 go build -o ${TASK_BINARY} ./api
	@echo "DONE!"