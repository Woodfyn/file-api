.PHONY:
.SILENT:

init:
	go mod download

build: init
	docker-compose build

run:
	docker-compose up