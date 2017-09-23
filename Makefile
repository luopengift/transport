APP = transport
MAIN = cmd/main.go
PREFIX = /usr/local

build:
	go build -o $(APP) $(MAIN)
update:
	go get -u ./...
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
.PHONY: build update install fmt lint test test-all clean all 
