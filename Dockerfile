FROM golang:1.23.8-alpine

RUN apk add tzdata git

WORKDIR /go/src/docker-image-updater
COPY . .

RUN go mod vendor
RUN go build ./...
RUN go install -v ./...

ENTRYPOINT ["/go/bin/docker-image-updater"]
CMD ["-update", "debian:latest", "-interval", "10"]