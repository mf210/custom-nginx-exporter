FROM golang:1.21.5 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY cmd/exporter/ ./cmd/exporter
COPY basic_collector.go basic.go ./

RUN go build -o /exporter cmd/exporter/main.go

FROM debian:latest

COPY --from=build /exporter /exporter

ENTRYPOINT ["/exporter"]