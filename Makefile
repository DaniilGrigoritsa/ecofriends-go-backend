SPOURCE_PATH = ./server.go
BIN_FILE = ./bin/server

format:
	go fmt ./...

tidy:
	go mod tidy

build:
	go build -o ${BIN_FILE} ${SPOURCE_PATH}

run:
	${BIN_FILE}

start: build run

generate-swag:
	swag init -g ${SPOURCE_PATH} -o ./docs

.PHONY: format tidy