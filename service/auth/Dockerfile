FROM golang:1.22.4-alpine as builder

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /opt

COPY go.mod go.sum ./
COPY pkg ./pkg
COPY service/auth ./service/auth

RUN go mod download

RUN go build -o app/main service/auth/cmd/app/main.go


FROM scratch as runtime

WORKDIR /opt/service/auth

COPY --from=builder /opt/app/main /main
ENTRYPOINT ["/main"]
