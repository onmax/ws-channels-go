# Makefile for Channels using WebSocket in Golang

include .env

help:
	@echo ""
	@echo "usage: make COMMAND"
	@echo ""
	@echo "Commands:"
	@echo "  run                 Run the application"
	@echo "  test                Run tests of the application"
	@echo "  docker-run          Run the application using Docker"
	@echo "  docker-test         Run tests of the application using Docker"
	@echo "  docker-stop         Stop the application"
	@echo "  docker-remove       Remove all the containers and images of the application"
	@echo "  docker-remove-test  Remove the image of the docker-test"
	@echo "  documentation       Open the docs in your default browser"

docker-run:
	@docker build -t $(CHANNELS-WS-DOCKER-NAME) . -f ./src/Dockerfile
	@docker run -d -p $(PORT):$(PORT) --name $(CHANNELS-WS-DOCKER-NAME) $(CHANNELS-WS-DOCKER-NAME)

docker-stop:
	@docker container stop -t 0 $(CHANNELS-WS-DOCKER-NAME)

docker-test:
	@docker build -t $(CHANNELS-WS-DOCKER-NAME-TEST) . -f tests/Dockerfile

docker-remove:
	@make docker-stop
	@docker container rm $(CHANNELS-WS-DOCKER-NAME)
	@docker image rm $(CHANNELS-WS-DOCKER-NAME)

docker-remove-test:
	@docker image rm $(CHANNELS-WS-DOCKER-NAME-TEST)

documentation:
	@xdg-open ./docs/index.html