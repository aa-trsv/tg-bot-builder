.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/main.go

run: build
	./.bin/bot

build-image:
	 docker build -t telegram-bot-otrs-builder:v1.0 .

start-container:
	docker run --name tg-bot-otrs-builder -p 8088:8088 --mount source=/opt/otrs,target=/opt/otrs telegram-bot-otrs-builder:v1.0