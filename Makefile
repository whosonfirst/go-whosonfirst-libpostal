export GOPATH=${PWD}

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src; then rm -rf src; fi
	cp -r vendor src
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-libpostal
	cp -r http src/github.com/whosonfirst/go-whosonfirst-libpostal/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:	rmdeps
	@GOPATH=$(GOPATH) go get -u "github.com/facebookgo/grace/gracehttp"
	@GOPATH=$(GOPATH) go get -u "github.com/openvenues/gopostal/expand"
	@GOPATH=$(GOPATH) go get -u "github.com/openvenues/gopostal/parser"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-sanitize"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-log"

vendor-deps: deps
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt http/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-libpostal-server cmd/wof-libpostal-server.go

data:
	bunzip2 -k testdata/*.bz2
