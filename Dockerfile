FROM golang:1.26-alpine

RUN apk add tzdata git

WORKDIR /go/src/docker-image-updater

COPY src src
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod vendor
RUN go build ./...
RUN go install -v ./...

ENTRYPOINT ["/go/bin/docker-image-updater"]
CMD ["-update", "debian:latest", "-interval", "10"]