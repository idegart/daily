ifneq (,$(wildcard ./.env))
    include .env
    export
endif

bold := $(shell tput bold)
sgr0 := $(shell tput sgr0)

start: start-db start-bot

stop: stop-db stop-bot

bot: build-bot run-bot

start-db:
	docker-compose up -d --remove-orphans db

stop-db:
	docker-compose stop db

migration:
	migrate create -ext sql -dir migrations ${name}

db-migrate:
	@migrate -path migrations -database "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable" up

start-bot:
	docker-compose up -d --build --remove-orphans bot

stop-bot:
	docker-compose stop bot

build-bot:
	go build -v -o bin/bot ./cmd/bot

run-bot:
	./bin/bot