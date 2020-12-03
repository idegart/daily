ifneq (,$(wildcard ./.env))
    include .env
    export
endif

bold := $(shell tput bold)
sgr0 := $(shell tput sgr0)

start: start-db start-apiserver start-dailybot

stop: stop-db stop-apiserver stop-dailybot

start-db:
	docker-compose up -d --remove-orphans db

stop-db:
	docker-compose stop db

migration:
	@migrate create -ext sql -dir migrations ${name}

db-migrate:
	@migrate -path migrations -database "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable" up

start-apiserver:
	docker-compose up -d --build --remove-orphans apiserver

stop-apiserver:
	docker-compose stop apiserver

start-dailybot:
	docker-compose up -d --build --remove-orphans dailybot

stop-dailybot:
	docker-compose stop dailybot

build-apiserver:
	go build -v -o bin/apiserver ./cmd/apiserver

build-dailybot:
	go build -v -o bin/dailybot ./cmd/dailybot