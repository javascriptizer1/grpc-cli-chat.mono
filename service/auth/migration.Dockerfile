FROM alpine:3.20

RUN apk update && apk upgrade && apk add bash && rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.20.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /opt

COPY service/auth/migrations ./migrations

ENTRYPOINT ["/bin/sh", "-c", "/bin/goose -dir ./migrations/postgres postgres up"]