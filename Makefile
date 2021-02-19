#!make
include .env

daily:
	go run ./cmd/daily

migration:
	migrate create -ext sql -dir migrations ${name}