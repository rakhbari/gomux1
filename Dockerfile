# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY config.yaml ./

RUN go build .

EXPOSE 8080
EXPOSE 8443

ENTRYPOINT [ "./gomux1" ]
