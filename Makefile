DB_DOCKER_CONTAINER=event_db
BINARY_NAME=event

postgres:
	docker run --name $(DB_DOCKER_CONTAINER) -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:alpine

createdb:
	docker exec -it $(DB_DOCKER_CONTAINER) createdb --username=root --owner=root event
run: 
	go run cmd/main.go
