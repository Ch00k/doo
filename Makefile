SHELL := /usr/bin/env bash

export DOO_DB_HOST ?= localhost
export DOO_DB_PORT ?= 5432
export DOO_DB_USER ?= doo
export DOO_DB_PASSWORD ?= doo
export DOO_DB_NAME ?= doo

DOCKER_IMAGE_NAME := ch00k/doo


build:
	go build

test:
	go test -v

run: build
	./doo

build_docker_image: build
	docker build -t ${DOCKER_IMAGE_NAME}:latest .

startdb:
	docker-compose up db

stopdb:
	docker-compose down db
