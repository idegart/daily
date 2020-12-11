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
	cd docker && docker-compose up -d --remove-orphans db

stop-db:
	cd docker && docker-compose stop db

migration:
	migrate create -ext sql -dir migrations ${name}

migrate:
	@migrate -path migrations -database "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable" up

start-bot:
	 cd docker && docker-compose up --build --remove-orphans bot

stop-bot:
	cd docker && docker-compose stop bot

build-bot:
	go build -v -o bin/bot ./cmd/bot

run-bot:
	./bin/bot