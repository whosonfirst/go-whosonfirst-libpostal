CWD=$(shell pwd)
GOPATH := $(CWD)

build:	fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src/github.com/whosonfirst/go-whosonfirst-geojson; then rm -rf src/github.com/whosonfirst/go-whosonfirst-geojson; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-geojson
	cp geojson.go src/github.com/whosonfirst/go-whosonfirst-geojson/geojson.go
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:   rmdeps
	@GOPATH=$(GOPATH) go get -u "github.com/jeffail/gabs"
	@GOPATH=$(GOPATH) go get -u "github.com/dhconnelly/rtreego"
	@GOPATH=$(GOPATH) go get -u "github.com/kellydunn/golang-geo"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-crawl"

vendor-deps: deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt *.go

bin:	self
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-contains cmd/wof-geojson-contains.go
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-dump cmd/wof-geojson-dump.go
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-enspatialize cmd/wof-geojson-enspatialize.go
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-polygons cmd/wof-geojson-polygons.go
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-validate cmd/wof-geojson-validate.go
