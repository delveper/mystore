.PHONY:run bench

ENV := .env
include $(ENV)

# App
run:
	go run ./cmd/main.go


# Docker
DOCKER_CONFIG_FLAGS := --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV) --log-level $(LOG_LEVEL)

docker-build:
	docker-compose ${DOCKER_CONFIG_FLAGS} build
docker-logs:
	docker-compose logs -f
docker-build-verbose:
	docker-compose --verbose ${DOCKER_CONFIG_FLAGS} build
docker-rebuild: docker-clean
	docker-compose ${DOCKER_CONFIG_FLAGS} build --no-cache
docker-up:
	sudo docker-compose ${DOCKER_CONFIG_FLAGS} up --detach
docker-start:
	sudo docker-compose ${DOCKER_CONFIG_FLAGS} start
docker-down:
	docker-compose ${DOCKER_CONFIG_FLAGS} down
docker-clean:
	docker-compose ${DOCKER_CONFIG_FLAGS} down --remove-orphans
docker-stop:
	docker-compose ${DOCKER_CONFIG_FLAGS} stop
docker-restart:
	docker-compose ${DOCKER_CONFIG_FLAGS} stop
	docker-compose ${DOCKER_CONFIG_FLAGS} run --detach