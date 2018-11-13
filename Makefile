PREFIX = /usr/local
GEOIP = utils/GeoLite2-City.mmdb
MAIN = cmd/*.go
APP = transport
VERSION	= 1.0.0
TIME = $(shell date "+%F %T")
GIT = $(shell git rev-parse HEAD)
PKG = github.com/luopengift/version

FLAG = "-X '${PKG}.VERSION=${VERSION}' -X '${PKG}.APP=${APP}' -X '${PKG}.TIME=${TIME}' -X '${PKG}.GIT=${GIT}'"

build: 
	go build -ldflags $(FLAG) -o ${APP} ${MAIN}
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
