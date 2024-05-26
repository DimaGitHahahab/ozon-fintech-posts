# Running app with postgres image
compose-up:
	docker compose up --build && docker compose logs --follow
compose-down:
	docker compose down --remove-orphans

# Running only app
APP_NAME = posts-app
DOCKER_IMAGE = $(APP_NAME):latest
DOCKER_CONTAINER = $(APP_NAME)-container

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run -p 8080:8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

clean:
	docker rm -f $(DOCKER_CONTAINER) || true
	docker rmi $(DOCKER_IMAGE) || true