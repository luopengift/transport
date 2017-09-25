APP = transport
MAIN = cmd/main.go
PREFIX = /usr/local
GEOIP = utils/GeoLite2-City.mmdb
build:
	go build -o $(APP) $(MAIN)
update:
	go get -u ./...
package: build
	tar -cvf $(APP).tar.gz $(APP) config.json init.sh Makefile
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
.PHONY: build update package install fmt lint test test-all clean all 
