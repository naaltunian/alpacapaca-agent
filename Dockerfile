FROM golang:1.16-alpine as build

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o paca-agent

FROM ubuntu:18.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates

COPY --from=build /build/paca-agent /usr/local/bin

CMD [ "paca-agent" ]