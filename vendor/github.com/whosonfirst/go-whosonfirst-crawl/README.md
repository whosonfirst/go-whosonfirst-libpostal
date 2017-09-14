# go-mapzen-whosonfirst-crawl

Go tools and libraries for crawling a directory of Who's On First data

## Usage

_Please rewrite me..._

## To do

* Documentation
* Proper error handling
* Remove GeoJSON specific stuff (or at least move it in to its own little playground)

## Caveats

This package relies on [a fork of the origin walk package](https://github.com/whosonfirst/walk) that relies on `runtime.GOMAXPROCS` to determine the number of concurrent processes used to crawl a directory tree.

## Tools

### wof-crawl-validate

For example:

```
./bin/wof-crawl-validate /usr/local/data/whosonfirst-data/data
validate JSON files in /usr/local/data/whosonfirst-data/data
2017/04/08 00:32:31 failed to parse /usr/local/data/whosonfirst-data/data/858/660/49/85866049.geojson, because invalid character '<' looking for beginning of object key string
2017/04/08 00:32:58 walked 504915 files (and 0 dirs) in 85.370 seconds
2017/04/08 00:32:58 okay 504914 errors 1
```

_Note: This only validates JSON-iness and not WOF-iness. Maybe someday it will do the latter but today it does not._

## See also

* https://github.com/whosonfirst/walk