package main

import (
	"flag"
	"fmt"
	geojson "github.com/whosonfirst/go-whosonfirst-geojson"
)

func main() {

	/*

			# ASSUMING
		        # lat := 45.523668
		        # lon := -73.600159

			# PLAIN VANILLA POLYGON

			$> ./bin/pip /usr/local/mapzen/whosonfirst-data/data/404/529/181/404529181.geojson /usr/local/mapzen/whosonfirst-data/data/857/848/31/85784831.geojson
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

			$> ./bin/pip /usr/local/mapzen/whosonfirst-data/data/136/251/273/136251273.geojson
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
			---

	*/

	flag.Parse()
	args := flag.Args()

	lat := 45.523668
	lon := -73.600159

	for _, path := range args {

		f, parse_err := geojson.UnmarshalFile(path)

		if parse_err != nil {
			panic(parse_err)
		}

		polygons := f.GeomToPolygons()
		contains := false

		fmt.Printf("%s has this many polygons: %d\n", path, len(polygons))

		for i, poly := range polygons {

			fmt.Printf("%s #%d has %d interior rings and a total of %d points\n", path, (i + 1), len(poly.InteriorRings), poly.CountPoints())

			c := poly.Contains(lat, lon)

			fmt.Printf("%s #%d contains point %t\n", path, (i + 1), c)

			if c {
				contains = true
			}
		}

		fmt.Printf("%s contains point: %t\n", path, contains)

		fmt.Printf("%s f.Contains() point: %t\n", path, f.Contains(lat, lon))
		fmt.Println("---")
	}

}
