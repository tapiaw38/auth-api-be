FROM golang:1.24-alpine as builder

LABEL maintainer="tapiaw38 Singh <tapiaw38@gmail.com>"

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

RUN apk update && apk add --no-cache postgresql-client curl && \
    curl -L -o /usr/local/bin/migrate https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz && \
    chmod +x /usr/local/bin/migrate

RUN apk add --no-cache build-base

COPY . /app/

EXPOSE 8082

COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]