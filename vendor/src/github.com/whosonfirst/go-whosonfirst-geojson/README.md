# go-whosonfirst-geojson

Go tools for working with Who's On First documents

## Usage

```
package main

import (
	"flag"
	"fmt"
	geojson "github.com/whosonfirst/go-whosonfirst-geojson/whosonfirst"
)

func main() {

	flag.Parse()
	args := flag.Args()

	for _, path := range args {

		// This is mostly just helper code to read the file
		// and call geojson.UnmarshalFeature (from a bag of bytes)

		f, parse_err := geojson.UnmarshalFile(path)

		if parse_err != nil {
			panic(parse_err)
		}

		fmt.Printf("# %s\n", path)
		fmt.Println(f.Dumps())
	}

}
```

## The longer version

This isn't really a "GeoJSON" specific library, yet. Right now it's just a thin wrapper around the [Gabs](https://github.com/jeffail/gabs) utility for wrangling unknown JSON structures in to a Go `WOFFeature` struct.

Eventually it would be nice to make Gabs hold hands with Paul Mach's [go.geojson](https://github.com/paulmach/go.geojson) and use the former to handle the GeoJSON properties dictionary. But that day is not today.

## The longer longer version

Right now this library has evolved and grown functionality on as-needed basis, targeting on Who's On First specific use-cases. As such it consists of a handful of WOF struct types - `WOFFeature` and `WOFPolygon` and `WOFSpatial` - that are wrappers around other people's heavy-lifting. There are not any WOF related interfaces but that's really the direction we want to head in... but we're not there yet. So things will probably change in the short-term. Not too much , hopefully.

## Utilities

Things you can find in the `cmd` and ultimately the `bin` directories.

### wof-geojson-contains

A tool for testing wether a given latitude and longitude is contained by one or more GeoJSON files. _As of this writing this tool lacks command line parameters for defining latitide and longitude._

```
# ASSUMING
# lat := 45.523668
# lon := -73.600159

# PLAIN VANILLA POLYGON

$> ./bin/wof-geojson-contains /usr/local/mapzen/whosonfirst-data/data/404/529/181/404529181.geojson /usr/local/mapzen/whosonfirst-data/data/857/848/31/85784831.geojson
/usr/local/mapzen/whosonfirst-data/data/404/529/181/404529181.geojson has this many polygons: 1
/usr/local/mapzen/whosonfirst-data/data/404/529/181/404529181.geojson #1 has 3 interior rings
/usr/local/mapzen/whosonfirst-data/data/404/529/181/404529181.geojson #1 contains point false
/usr/local/mapzen/whosonfirst-data/data/404/529/181/404529181.geojson contains point: false
/usr/local/mapzen/whosonfirst-data/data/404/529/181/404529181.geojson f.Contains() point: false
---
/usr/local/mapzen/whosonfirst-data/data/857/848/31/85784831.geojson has this many polygons: 1
/usr/local/mapzen/whosonfirst-data/data/857/848/31/85784831.geojson #1 has 0 interior rings
/usr/local/mapzen/whosonfirst-data/data/857/848/31/85784831.geojson #1 contains point true
/usr/local/mapzen/whosonfirst-data/data/857/848/31/85784831.geojson contains point: true
/usr/local/mapzen/whosonfirst-data/data/857/848/31/85784831.geojson f.Contains() point: true
---

# MULTI POLYGON

$> ./bin/wof-geojson-contains /usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson has this many polygons: 8
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #1 has 0 interior rings
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #1 contains point false
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #2 has 0 interior rings
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #2 contains point false
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #3 has 0 interior rings
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #3 contains point false
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #4 has 0 interior rings
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #4 contains point true
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #5 has 0 interior rings
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #5 contains point false
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #6 has 0 interior rings
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #6 contains point false
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #7 has 0 interior rings
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #7 contains point false
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #8 has 0 interior rings
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson #8 contains point false
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson contains point: true
/usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson f.Contains() point: true
```

### wof-geojson-dump

Print the ID, name and placetype for one or more GeoJSON files. This is a utility to test the `Id` and `Name` and `Placetype` methods for a GeoJSON document parsed by `go-whosonfirst-geojson`

```
$> ./bin/wof-geojson-dump /usr/local/mapzen/whosonfirst-data/data/101/736/545/101736545.geojson
# /usr/local/mapzen/whosonfirst-data/data/101/736/545/101736545.geojson
ID is 101736545
Name is Montréal
Placetype is locality
```

### wof-geojson-enspatialize

This is a utility for testing the `SpatializeGeom` functionality for one or more GeoJSON files.

```
./bin/wof-geojson-enspatialize /usr/local/mapzen/whosonfirst-data/data/101/736/545/101736545.geojson
Enspatialize bounding box
&{0xc210038d20 101736545 Montréal locality -1}
101736545
Enspatialize geom
1
&{0xc210038de0 101736545 Montréal locality 0}
```

### wof-geojson-polygons

This is a utility for testing the `GeomToPolygons` functionality, by printing the number of points in each outer ring, for one or more GeoJSON files.

```
./bin/wof-geojson-polygons /usr/local/mapzen/whosonfirst-data/data/101/736/545/101736545.geojson
5206 points
```

### wof-geojson-validate

Validate a directory full of GeoJSON files. Specifically validate that they are _valid JSON_ and nothing else. A more full-feature Who's On First validator in Go may be written in the future but today is not that day. You could take a look at [py-mapzen-whosonfirst-validator](https://github.com/whosonfirst/py-mapzen-whosonfirst-validator) for that.

```
$> /bin/wof-geojson-validate -source /usr/local/mapzen/whosonfirst-data/data/ -processes 200
time to validate 401499 files: 1m11.002168559s
```

## See also

* https://www.github.com/jeffail/gabs
