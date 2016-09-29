# go-whosonfirst-libpostal

Go tools for working with libpostal (sometimes in the service of Who's On First)

## Install

```
make bin
```

## wof-libpostal-server

```
$> wof-libpostal-server -options

Usage of wof-libpostal-server:
  -gracehttp.log
	Enable logging. (default true)
  -host string
    	The hostname to listen for requests on (default "localhost")
  -pidfile string
    	   Where to write a PID file for wof-libpostal-server. If empty the PID file will be written to wof-libpostal-server.pid in the current directory
  -port int
    	The port number to listen for requests on (default 8080)
```

_Note you will need to install the underlying [libpostal](https://github.com/openvenues/libpostal) C library yourself in order for `wof-libpostal-server` to work._

`wof-libpostal-server` exposes the following endpoints:

### Endpoints

#### GET /expand _?address=ADDRESS_

This endpoint accepts a single `address` parameter and expands it into one or more normalized forms suitable for geocoder queries.

```
curl -s -X GET 'http://localhost:8080/expand?address=475+Sansome+St+San+Francisco+CA' | python -mjson.tool
[
    "475 sansome saint san francisco california",
    "475 sansome saint san francisco ca",
    "475 sansome street san francisco california",
    "475 sansome street san francisco ca"
]
```

#### GET /parse _?address=ADDRESS_

This endpoint accepts a single `address` parameter and parses it in to its components.

```
curl -s -X GET 'http://localhost:8080/parse?address=475+Sansome+St+San+Francisco+CA' | python -mjson.tool
[
    {
        "label": "house_number",
        "value": "475"
    },
    {
        "label": "road",
        "value": "sansome st"
    },
    {
        "label": "city",
        "value": "san francisco"
    },
    {
        "label": "state",
        "value": "ca"
    }
]
```

#### GET /debug/vars

This endpoint exposes all the usual Go [expvar](https://golang.org/pkg/expvar/) debugging output along with the following additional properies:

* AvgParse
* AvgExpand
* ErrInput
* ErrRead
* ErrMarshal
* ReqExpand
* ReqParse
* SuccessExpand
* SuccessParse

_Note: This endpoint is only available from the machine the server is running on._

## wof-libpostal-crawl

```
$> wof-libpostal-crawl -options /path/to/whosonfirst-data-venue-us/data

Usage of ./bin/wof-libpostal-crawl:
  -libpostal-host string
    		  The host for the libpostal endpoint
  -libpostal-port int
    		  The host for the libpostal port (default 8080)
  -output string
    	  Where to write output data (default "libpostal.csv")
  -processes int
    	     The number of concurrent processes to clone data with (default 16)
  -properties string
    	      A comma-separated list of GeoJSON properties used to construct an address (default "sg:address,sg:city,sg:province,sg:postcode")
```

`wof-libpostal-crawl` traverses a directory containin GeoJSON files and queries an instance of `wof-libpostal-server` with an address string derived from one or more properties in the GeoJSON feature. For example, this:

```
./bin/wof-libpostal-crawl -libpostal-host localhost -processes 200 /usr/local/data/whosonfirst-data-venue-au/data/
parsed 697614 files in /usr/local/data/whosonfirst-data-venue-au/data/ in 314.434 seconds avg ttq: 836 ms
```

Produces a file that looks like this:

```
id,address,results
336966001,343 Montague Rd West End QLD 4810,"[{""label"":""house_number"",""value"":""343""},{""label"":""road"",""value"":""montague rd""},{""label"":""city_district"",""value"":""west end""},{""label"":""state"",""value"":""qld""},{""label"":""postcode"",""value"":""4810""}]"
186121023,3/ 53 Smith St Kempsey NSW 2440,"[{""label"":""level"",""value"":""3""},{""label"":""house_number"",""value"":""/ 53""},{""label"":""road"",""value"":""smith st""},{""label"":""suburb"",""value"":""kempsey""},{""label"":""state"",""value"":""nsw""},{""label"":""postcode"",""value"":""2440""}]"

... and so on
```

## Performance and load testing

This assumes the server a pair of single-CPU m3.medium instances, fronted by a load-balancer, and a URLs file containing 1517965 addresses in California. There is also a URLs file with 18M addresses from most of SimpleGeo but that is 8GB and I got bored waiting for [siege](https://www.joedog.org/siege-home/) to it to load it in to memory...

### expand

```
$> siege -c 500 -i -f urls-expand.txt
** SIEGE 3.0.5
** Preparing 500 concurrent users for battle.
The server is now under siege...  C-c C-c
Lifting the server siege...      done.

Transactions:		      147245 hits
Availability:		      100.00 %
Elapsed time:		      148.14 secs
Data transferred:	       32.60 MB
Response time:		        0.00 secs
Transaction rate:	      993.96 trans/sec
Throughput:		        0.22 MB/sec
Concurrency:		        2.93
Successful transactions:      147072
Failed transactions:	           0
Longest transaction:	        1.16
Shortest transaction:	        0.00
```

### parse

```
$> siege -c 500 -i -f urls-parse.txt 
** SIEGE 3.0.5
** Preparing 500 concurrent users for battle.
The server is now under siege...  C-c C-c
Lifting the server siege...      done.

Transactions:		      131044 hits
Availability:		      100.00 %
Elapsed time:		      131.96 secs
Data transferred:	       23.17 MB
Response time:		        0.00 secs
Transaction rate:	      993.06 trans/sec
Throughput:		        0.18 MB/sec
Concurrency:		        3.85
Successful transactions:      130895
Failed transactions:	           0
Longest transaction:	        5.06
Shortest transaction:	        0.00 
```

## See also

* https://github.com/openvenues/libpostal
* https://github.com/openvenues/gopostal
