FROM golang:1.20-alpine3.18 AS builder

RUN apk update && \
    apk add --no-cache git && \
    apk add --no-cache build-base # it's necessary for sqlite.

ENV GOPATH="$HOME/go"
WORKDIR $GOPATH/src/github.com/carloseduribeiro/crud-go/

ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN GOOS=linux go mod download -x

# it's necessary for sqlite:
ENV CGO_ENABLED=1

COPY . .
RUN go build -C ./cmd/server -o /go/bin/crud-go

FROM alpine:3.18

COPY --from=builder /go/bin/crud-go /go/bin/crud-go
COPY .env /go/bin/.env

WORKDIR /go/bin

ENTRYPOINT ["/go/bin/crud-go"]