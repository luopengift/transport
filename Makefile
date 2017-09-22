APP = transport
MAIN = cmd/main.go
PREFIX = /usr/local
build:
	go build -o $(APP) $(MAIN)
install:
	mv -f $(APP) $(PREFIX)/bin
fmt:
	go fmt ./...
lint:
	go vet ./...
test:
	go test -short ./...
test-all: lint
	go test ./...
clean:
	rm -f $(APP)
	rm -f $(APP).exe
	rm -f $(PREFIX)/bin/$(APP)
.PHONY: build install fmt lint test test-all clean all 
