# go-whosonfirst-libpostal

Go tools for working with libpostal (sometimes in the service of Who's On First)

## Install

Make sure you have [`Go`](https://golang.org/doc/install) and [`libpostal`](https://github.com/openvenues/libpostal) installed. Then

```
make bin
```

## wof-libpostal-server

```
$> ./bin/wof-libpostal-server -options

Usage of wof-libpostal-server:
  -gracehttp.log
	Enable logging. (default true)
  -host string
    	The hostname to listen for requests on (default "localhost")
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

## See also

* https://github.com/openvenues/libpostal
* https://github.com/openvenues/gopostal
