FROM alpine:3.20

RUN apk update && apk upgrade && apk add bash && rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.20.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /opt

COPY service/auth/migrations ./migrations

ENV GOOSE_DBSTRING=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable
RUN /bin/goose -dir ./migrations/postgres postgres up