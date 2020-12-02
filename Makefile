ifneq (,$(wildcard ./.env))
    include .env
    export
endif

bold := $(shell tput bold)
sgr0 := $(shell tput sgr0)

start:
	@printf "$(bold)Starting database$(sgr0)\n"
	@make start-db
	@printf "$(bold)Database started\nStarting API Server$(sgr0)\n"
	@make start-apiserver

build-apiserver:
	go build -v -o bin/apiserver ./cmd/apiserver

start-db:
	docker-compose up -d --remove-orphans db

start-apiserver:
	docker-compose up --build --remove-orphans apiserver