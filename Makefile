help:
	@echo "Available commands:"
	@echo "\trun                - runs the exporter"
	@echo "\twatch              - runs the exporter with hot reload"
	@echo "\ttest               - runs the tests"
	@echo "\tpodman-build       - builds the podman container"
	@echo "\tpodman-run         - runs the podman container"
	@echo ""

.PHONY: run
run:
	go build && ./shelly-exporter \
		--config config.yaml \
		--web.listen-address :9123

.PHONY: watch
watch:
	go get -u github.com/cosmtrek/air
	air

.PHONY: test
test:
	go test ./... -timeout 30s -v -cover


.PHONY: podman-build
podman-build:
	podman build -t shelly-exporter .


.PHONY: podman-run
podman-run:
	podman run --name shelly-exporter \
		-v $(shell pwd)/config.yaml:/app/config.yaml \
		-p 127.0.0.1:9123:9123/tcp \
		--rm -it shelly-exporter:latest