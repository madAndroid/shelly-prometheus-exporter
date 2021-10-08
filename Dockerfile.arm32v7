FROM arm32v7/golang as dependencies
WORKDIR /go/src/app
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM dependencies as builder
COPY . .
RUN go test ./... -timeout 30s -cover
RUN CGO_ENABLED=0 go build -o shelly-exporter

FROM arm32v7/alpine
LABEL maintainer="Alex Voigt <mail@alexander-voigt.info>"
WORKDIR /app/
VOLUME ["/app"]
COPY --from=builder /go/src/app/shelly-exporter .
CMD ["./shelly-exporter"]