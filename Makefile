.PHONY:
.SILENT:

init:
	go mod download

build: init
	docker-compose build

run:
	docker-compose up

docker-clear:
	docker builder prune -f
	docker image prune -f
	docker container prune -f