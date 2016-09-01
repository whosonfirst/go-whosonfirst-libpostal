prep:
	if test -d pkg; then rm -rf pkg; fi

self:	prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-crawl; then rm -rf src/github.com/whosonfirst/go-whosonfirst-crawl; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-crawl
	cp crawl.go src/github.com/whosonfirst/go-whosonfirst-crawl/crawl.go

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	rmdeps deps fmt bin

deps:   self
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/walk"

fmt:
	go fmt cmd/*.go
	go fmt *.go

bin:	self
	@GOPATH=$(shell pwd) go build -o bin/wof-count cmd/wof-count.go
	@GOPATH=$(shell pwd) go build -o bin/wof-crawl-dtwt cmd/wof-crawl-dtwt.go
