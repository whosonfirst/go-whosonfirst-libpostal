CWD=$(shell pwd)
GOPATH := $(CWD)

build:	rmdeps deps fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-geojson; then rm -rf src/github.com/whosonfirst/go-whosonfirst-geojson; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-geojson
	cp geojson.go src/github.com/whosonfirst/go-whosonfirst-geojson/geojson.go

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:   self
	@GOPATH=$(GOPATH) go get -u "github.com/jeffail/gabs"
	@GOPATH=$(GOPATH) go get -u "github.com/dhconnelly/rtreego"
	@GOPATH=$(GOPATH) go get -u "github.com/kellydunn/golang-geo"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-crawl"

fmt:
	go fmt cmd/*.go
	go fmt *.go

bin:	self
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-contains cmd/wof-geojson-contains.go
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-dump cmd/wof-geojson-dump.go
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-enspatialize cmd/wof-geojson-enspatialize.go
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-polygons cmd/wof-geojson-polygons.go
	@GOPATH=$(GOPATH) go build -o bin/wof-geojson-validate cmd/wof-geojson-validate.go
