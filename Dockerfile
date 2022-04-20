FROM golang:1.18-alpine3.15 AS builder

COPY . /github.com/aa-trsv/telegram-bot-otrs-builder/
WORKDIR /github.com/aa-trsv/telegram-bot-otrs-builder/

RUN go mod download
RUN go build -o ./bin/bot cmd/main.go


FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/aa-trsv/telegram-bot-otrs-builder/bin/bot .

EXPOSE 8088

CMD ["./bot"]