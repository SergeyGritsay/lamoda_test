.PHONY: all

include .env
export

PROJ_PATH := ${CURDIR}
DOCKER_PATH := ${PROJ_PATH}/docker

APP=lamoda-test
MIGRATION_TOOL=goose
MIGRATION_DIR=./db/migrations

BASIC_IMAGE=dep

build:
	go build -o .bin/${APP} cmd/${APP}/main.go
	chmod ugo+x .bin/${APP}

build-docker:
	sudo rm -rf .database/
	docker build -t ${APP}-image -f ${DOCKER_PATH}/${APP}.dockerfile .

app-setup-and-up: build-docker app-up

app-up: build
	docker-compose up

all: app-setup-and-up

app-bash: 
	docker-compose run --rm --no-deps --name ${APP}-service warehouse ash

app-up-local: build	
	./.bin/lamoda-test

db-bash:
	docker-compose run --rm --no-deps --name ${APP}-db db ash

package-tidy:
	go mod tidy

db-up:
	docker-compose run --rm --no-deps --name ${APP}-db db ash

up_migrate:
	go install github.com/pressly/goose/v3/cmd/goose@latest

	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=qwerty dbname=warehouse host=0.0.0.0 sslmode=disable" goose -dir ./db/migrations up
down_migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=qwerty dbname=warehouse host=0.0.0.0 sslmode=disable" goose -dir ./db/migrations down
