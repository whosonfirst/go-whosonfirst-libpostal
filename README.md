# go-whosonfirst-libpostal

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

## wof-libpostal-crawl

```
$> wof-libpostal-crawl -options /path/to/whosonfirst-data-venue-us/data

Usage of wof-libpostal-crawl:
  -libpostal-host string
    		  The host for the libpostal endpoint
  -libpostal-port int
    		  The host for the libpostal port (default 8080)
  -output string
    	  Where to write output data (default "libpostal.csv")
  -processes int
    	     The number of concurrent processes to clone data with (default 16)
```

_Note that `wof-libpostal-crawl` is currently designed for use with Who's On First documents with specific SimpleGeo prefixed keys: `sg:address, sg:city, sg:province, sg:postcode`. It will eventually be adapted for other things._

## See also

* https://github.com/openvenues/libpostal
* https://github.com/openvenues/gopostal