APP = transport
MAIN = cmd/main.go

build:
	go build -o $(APP) $(MAIN)
fmt:
	go fmt ./...
lint:
	go vet ./...
test:
	go test -short ./...
test-all: lint
	go test ./...
clean:
	rm -f transport
	rm -f transport.exe
.PHONY: fmt deps transport transport.exe install test test-windows lint test-all \
	package clean docker-run docker-run-circle docker-kill
