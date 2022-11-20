FROM golang:1.19.2-bullseye as dependencies
WORKDIR /go/src/app
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download -json -x

FROM dependencies as builder
COPY . .
RUN go test ./... -timeout 30s -cover
RUN go build -o shelly-exporter

FROM debian:bullseye-slim
#FROM alpine:latest
LABEL maintainer="Alex Voigt <mail@alexander-voigt.info>"
WORKDIR /app/
VOLUME ["/app"]
COPY --from=builder /go/src/app/shelly-exporter .
CMD ["./shelly-exporter"]
