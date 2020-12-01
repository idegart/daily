ifneq (,$(wildcard ./.env))
    include .env
    export
endif

bold := $(shell tput bold)
sgr0 := $(shell tput sgr0)

start:
	@printf "$(bold)Starting database$(sgr0)\n"
	@make start-db
	@printf "$(bold)Database started\nStarting app$(sgr0)\n"
	@make start-app

build:
	go build -v -o bin/${APP_NAME} ./cmd/bot

start-db:
	docker-compose up -d db

start-app:
	docker-compose up --build app