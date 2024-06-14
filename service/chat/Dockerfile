FROM golang:1.22.4-alpine as builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /opt

COPY go.mod go.sum ./
COPY pkg ./pkg
COPY service/chat ./service/chat

RUN go mod download

RUN go build -o app/main service/chat/cmd/app/main.go


FROM scratch as runtime

WORKDIR /opt/service/chat

COPY --from=builder /opt/app/main /main
ENTRYPOINT ["/main"]
